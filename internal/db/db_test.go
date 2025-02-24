package db

import (
	"errors"
	"simple-file-processor/internal/mocks/mockdb"
	"simple-file-processor/internal/models"
	"testing"

	"github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

var (
	l = zerolog.Nop()
)

func Test_NewDB_WhenCalled_ReturnsDB(t *testing.T) {
	db := new(mockdb.GormDB)
	g := gomega.NewWithT(t)
	gdb := NewDB(db, l)
	g.Expect(gdb).NotTo(gomega.BeNil())
}

func Test_Migrate_WhenCalled_ReturnsNil(t *testing.T) {
	db := new(mockdb.GormDB)
	g := gomega.NewWithT(t)
	gdb := NewDB(db, l)
	db.On("AutoMigrate", &models.File{}).Return(nil)
	err := gdb.Migrate()
	g.Expect(err).To(gomega.BeNil())
}

func Test_Migrate_WhenErrorAutoMigrate_ReturnsError(t *testing.T) {
	db := new(mockdb.GormDB)
	g := gomega.NewWithT(t)
	gdb := NewDB(db, l)
	db.On("AutoMigrate", &models.File{}).Return(errors.New("error"))
	err := gdb.Migrate()
	g.Expect(err).NotTo(gomega.BeNil())
}

func Test_InsertFileMetadata_WhenNoError_ReturnsNil(t *testing.T) {
	db := new(mockdb.GormDB)
	g := gomega.NewWithT(t)
	gdb := NewDB(db, l)
	file := &models.File{
		OriginalName: "test.txt",
	}

	db.On("Create", file).Return(&gorm.DB{Error: nil})
	err := gdb.InsertFileMetadata(file)
	g.Expect(err).To(gomega.BeNil())
}

func Test_InsertFileMetadata_WhenError_ReturnsError(t *testing.T) {
	db := new(mockdb.GormDB)
	g := gomega.NewWithT(t)
	gdb := NewDB(db, l)
	file := &models.File{
		OriginalName: "test.txt",
	}

	db.On("Create", file).Return(&gorm.DB{Error: errors.New("error")})
	err := gdb.InsertFileMetadata(file)
	g.Expect(err).NotTo(gomega.BeNil())
}
