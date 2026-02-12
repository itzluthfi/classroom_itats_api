package main_test

import (
	"classroom_itats_api/repositories"
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var _ = Describe("Main", Ordered, func() {
	var db *gorm.DB
	var user_repository repositories.UserRepository

	dsn := "host=localhost user=postgres password=123 dbname=ak070203 port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	Expect(err).ShouldNot(HaveOccurred())

	db = conn
	// db_connection.AutoMigrate(&entities.User{})

	user_repository = repositories.NewUserRepository(db)

	ctx := context.Background()

	Describe("User Repository", func() {
		Describe("Get User From Database", func() {
			When("Get One User From Database", func() {
				It("Should Return One User From Database", func() {
					user, err := user_repository.GetUserByUid(ctx, 43703)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(user.Name).To(Equal("06.2020.1.07351"))
					Expect(user.Mail).To(Equal("satyagray@gmail.com"))
				})
			})
		})
	})
})
