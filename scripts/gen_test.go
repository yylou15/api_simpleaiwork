package main

import (
	"log"
	"os"
	"strings"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gen"
	"gorm.io/gorm"
)

func TestGen(t *testing.T) {
	log.Println("Starting TestGen with SQLite")
	// Use file-based SQLite
	db, err := gorm.Open(sqlite.Open("gen.db"), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to connect to database:", err)
	}
	defer os.Remove("gen.db") // Clean up

	// Read SQL file
	sqlContent, err := os.ReadFile("init_sqlite.sql")
	if err != nil {
		t.Fatal("Failed to read init.sql:", err)
	}

	// Execute SQL
	log.Println("Executing init.sql...")
	queries := splitSQL(string(sqlContent))
	for _, q := range queries {
		if q == "" {
			continue
		}
		if err := db.Exec(q).Error; err != nil {
			t.Fatalf("Failed to execute query: %s\nError: %v", q, err)
		}
	}

	log.Println("Generating GORM code...")
	// Generate GORM code
	g := gen.NewGenerator(gen.Config{
		OutPath:      "../biz/say_right/dal/query",
		ModelPkgPath: "../biz/say_right/dal/model",
		Mode:         gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.UseDB(db)

	// Map SQLite types to Go types
	// Note: GORM Gen handles mapping, but SQLite might need hints for integers to be int64
	dataMap := map[string]func(gorm.ColumnType) (dataType string){
		"integer": func(detail gorm.ColumnType) (dataType string) {
			return "int64"
		},
		"text": func(detail gorm.ColumnType) (dataType string) {
			return "string"
		},
		"blob": func(detail gorm.ColumnType) (dataType string) {
			return "[]byte"
		},
		"datetime": func(detail gorm.ColumnType) (dataType string) {
			return "time.Time"
		},
	}
	g.WithDataTypeMap(dataMap)

	// Generate all tables
	g.ApplyBasic(
		g.GenerateModel("users"),
		g.GenerateModel("user_identities"),
		g.GenerateModel("email_verifications"),
		g.GenerateModel("categories"),
		g.GenerateModel("templates"),
		g.GenerateModel("template_details"),
	)

	g.Execute()
	log.Println("Done.")
}

func splitSQL(sql string) []string {
	var stmts []string
	parts := strings.Split(sql, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			stmts = append(stmts, p)
		}
	}
	return stmts
}
