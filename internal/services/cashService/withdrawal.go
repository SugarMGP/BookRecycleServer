package cashService

import (
	"bookrecycle-server/internal/models"
	"bookrecycle-server/internal/services/userService"
	"bookrecycle-server/pkg/database"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func Withdrawal(userID uint, amount float64) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		user, err := userService.GetUserByID(userID)
		if err != nil {
			return err
		}
		balance, _ := decimal.NewFromString(user.Balance)
		cash := decimal.NewFromFloat(amount)
		user.Balance = balance.Sub(cash).StringFixedBank(2)
		if err = tx.Save(user).Error; err != nil {
			return err
		}

		return tx.Create(&models.Withdrawal{
			Money:   cash.StringFixedBank(2),
			Account: user.Phone,
			Name:    user.Name,
		}).Error
	})
}

func GetWithdrawalList() (withdrawalList []*models.Withdrawal, err error) {
	err = database.DB.Find(&withdrawalList).Error
	return
}
