package main

import (
	"fmt"
	"github.com/Baraulia/goLab/IndTask.git"
	"github.com/Baraulia/goLab/IndTask.git/internal/config"
	"github.com/Baraulia/goLab/IndTask.git/internal/repository"
	"github.com/Baraulia/goLab/IndTask.git/pkg/logging"
	"net/smtp"
	"strings"
	"time"
)

var debtors = make(chan []*IndTask.Debtor, 1)

func CheckData(ticker *time.Ticker, Rep *repository.Repository, logger logging.Logger) {
	listDebtors, err := CheckReturnData(Rep)
	if err != nil {
		logger.Errorf("Can not check return data for acts (%s):%s", time.Now(), err)
	}
	if len(listDebtors) > 0 {
		debtors <- listDebtors
	}
	for {
		select {
		case <-ticker.C:
			listDebtors, err := CheckReturnData(Rep)
			if err != nil {
				logger.Errorf("Can not check return data for acts (%s):%s", time.Now(), err)
			}
			if len(listDebtors) > 0 {
				debtors <- listDebtors
			}
		}
	}
}

func SendEmail(cfg *config.Config, logger logging.Logger, auth smtp.Auth) {
	for {
		select {
		case <-debtors:
			listDebtors := <-debtors
			from := cfg.Mail.From
			smtpHost := cfg.Mail.SmtpHost
			smtpPort := cfg.Mail.SmtpPort
			for _, debtor := range listDebtors {
				msg := fmt.Sprintf(" Уважаемый %s, Вам необходимо вернуть книги %s!", debtor.Name, debtor.Book)
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
}

func CheckReturnData(rep *repository.Repository) ([]*IndTask.Debtor, error) {
	listDeptors, err := rep.AppAct.CheckReturnData()
	if err != nil {
		return nil, fmt.Errorf("checkReturnData:%w", err)
	}
	return listDeptors, nil
}
