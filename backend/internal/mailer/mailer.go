package mailer

import (
	"log/slog"

	"gopkg.in/gomail.v2"

	"github.com/anujgupta/level-up-backend/internal/config"
)

type EmailJob struct {
	To       string
	Subject  string
	Template string
	Data     any
}

type Mailer struct {
	jobChan chan EmailJob
	dialer  *gomail.Dialer
	from    string
	logger  *slog.Logger
}

func New(cfg *config.Config, logger *slog.Logger) *Mailer {
	dialer := gomail.NewDialer(cfg.SMTPHost, cfg.SMTPPort, cfg.SMTPUser, cfg.SMTPPassword)
	return &Mailer{
		jobChan: make(chan EmailJob, 100),
		dialer:  dialer,
		from:    cfg.EmailFrom,
		logger:  logger,
	}
}

func (m *Mailer) Start(workers int) {
	for i := 0; i < workers; i++ {
		go m.worker()
	}
}

func (m *Mailer) Send(job EmailJob) {
	select {
	case m.jobChan <- job:
	default:
		m.logger.Warn("mailer queue full, dropping email", "to", job.To, "subject", job.Subject)
	}
}

// Close drains the job channel. Call after http.Server.Shutdown().
func (m *Mailer) Close() {
	close(m.jobChan)
}

func (m *Mailer) worker() {
	for job := range m.jobChan {
		if err := m.send(job); err != nil {
			m.logger.Error("failed to send email",
				"to", job.To,
				"subject", job.Subject,
				"err", err,
			)
		}
	}
}

func (m *Mailer) send(job EmailJob) error {
	body := renderTemplate(job.Template, job.Data)

	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", job.To)
	msg.SetHeader("Subject", job.Subject)
	msg.SetBody("text/html", body)

	return m.dialer.DialAndSend(msg)
}
