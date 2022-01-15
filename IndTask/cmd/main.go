package main

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/handler"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/internal/service"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"github.com/Baraulia/goLab/IndTask.git/pkg/postgres"
	"github.com/Baraulia/goLab/IndTask.git/pkg/server"
	_ "github.com/lib/pq"
	"net/smtp"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	logger := logging.GetLogger()
	cfg := config.GetConfig()

	db, err := postgres.NewPostgresDB(postgres.PostgresDB{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		logger.Panicf("failed to initialize db:%s", err.Error())
	}

	Rep := repository.NewRepository(db)
	services := service.NewService(Rep, cfg)
	handlers := handler.NewHandler(logger, services)
	srv := new(server.Server)
	logger.Infof("Running server on %s:%s...", cfg.Listen.BindIp, cfg.Listen.Port)
	go func() {
		if err := srv.Run(cfg.Listen.BindIp, cfg.Listen.Port, handlers.InitRoutes()); err != nil {
			logger.Panicf("Error occured while running http server: %s", err.Error())
		}
	}()

	logger.Info("IndTask started")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	ticker := time.NewTicker(24 * time.Hour)
	debtors := make(chan []IndTask.Debtor, 1)

	go func() {
		for {
			select {
			case <-ticker.C:
				listDebtors, err := CheckReturnData(Rep)
				if err != nil {
					logger.Errorf("Can not check return data for issue acts (%s):%s", time.Now(), err)
				}
				if len(listDebtors) > 0 {
					debtors <- listDebtors
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case <-debtors:
				listDebtors := <-debtors
				from := cfg.Mail.From
				password := cfg.Mail.Password
				smtpHost := cfg.Mail.SmtpHost
				smtpPort := cfg.Mail.SmtpPort
				auth := smtp.PlainAuth("", from, password, smtpHost)
				for _, debtor := range listDebtors {
					msg := fmt.Sprintf(" Уважаемый %s, Вам необходимо вернуть книгу %s!", debtor.Name, debtor.Book)
					message := strings.Replace("From: "+from+"~To: "+debtor.Email+"~Subject: "+cfg.Mail.Subject+"~~", "~", "\r\n", -1) + msg
					err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{debtor.Email}, []byte(message))
					if err != nil {
						logger.Errorf("Error while sending email to %s:%s", debtor.Email, err)
						return
					}
				}
				logger.Info("Email Sent Successfully!")
			}
		}

	}()

	<-quit
	ticker.Stop()

}
func CheckReturnData(rep *repository.Repository) ([]IndTask.Debtor, error) {
	listDeptors, err := rep.AppMove.CheckReturnData()
	if err != nil {
		return nil, err
	}
	return listDeptors, nil
}
