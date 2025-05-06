package database

import (
	modelOrder "github.com/andrianprasetya/eventHub/internal/Order/model"
	modelAuditSecurity "github.com/andrianprasetya/eventHub/internal/audit_security_log/model"
	modelCheckIn "github.com/andrianprasetya/eventHub/internal/check_in/model"
	modelCommunication "github.com/andrianprasetya/eventHub/internal/communication/model"
	modelEvent "github.com/andrianprasetya/eventHub/internal/event/model"
	modelReport "github.com/andrianprasetya/eventHub/internal/report/model"
	modelTenant "github.com/andrianprasetya/eventHub/internal/tenant/model"
	modelTicket "github.com/andrianprasetya/eventHub/internal/ticket/model"
	modelUser "github.com/andrianprasetya/eventHub/internal/user/model"
	"log"
)

// MigrateDatabase runs database migrations
func MigrateDatabase() {
	db := GetConnection()

	err := db.AutoMigrate(
		//tenant
		&modelTenant.Tenant{},
		&modelTenant.TenantSetting{},
		&modelTenant.SubscriptionPlan{},
		&modelTenant.Subscription{},

		//user
		&modelUser.Role{},
		&modelUser.User{},
		&modelUser.Permission{},

		//Event
		&modelEvent.Event{},
		&modelEvent.EventCategory{},
		&modelEvent.EventSession{},
		&modelEvent.EventTag{},
		&modelEvent.EventCustomField{},


		//Ticket
		&modelTicket.EventTicket{},
		&modelTicket.Discount{},
		&modelTicket.TicketCoupon{},
		&modelTicket.EventCustomField{},

		//Order
		&modelOrder.Order{},
		&modelOrder.OrderItem{},
		&modelOrder.Invoices{},
		&modelOrder.PaymentTransaction{},

		//CheckIn
		&modelCheckIn.Ticket{},
		&modelCheckIn.CheckIn{},

		//Communication
		&modelCommunication.Notification{},
		&modelCommunication.EmailTemplate{},
		&modelCommunication.Webhook{},

		//Report
		&modelReport.Report{},
		&modelReport.Export{},

		//Audit Security Log
		&modelAuditSecurity.ActivityLog{},
		&modelAuditSecurity.LoginHistory{},
		&modelAuditSecurity.ApiKey{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}
	log.Println("Migration completed successfully")
}
