# Dokumentasi Mendalam & Panduan Lengkap REST API Card Bridge (`bab6`)

Dokumentasi ini dibuat untuk membedah proyek Go secara total. Penjelasan disusun baris demi baris, fungsi demi fungsi, serta simbol teknis agar pemula memahami logika di balik aplikasi ini.

---

## 💡 1. Konsep Dasar & Analogi (Restoran)

Aplikasi ini menggunakan **Layered Architecture**. Bayangkan Restoran:

- **ROUTER (Kasir)**: Menerima pesanan dan mengarahkan ke Pelayan yang tepat.
- **CONTROLLER (Pelayan)**: Memeriksa format pesanan & menyajikan makanan (format JSON).
- **REPOSITORY (Koki)**: Mengolah data langsung di dapur (Query SQL ke Database).
- **DATABASE (Gudang)**: Tempat menyimpan data (PostgreSQL).

---

## 📖 2. Glosarium Simbol & Istilah Teknis

| Simbol / Istilah | Penjelasan Pemula |
| :--- | :--- |
| **`//go:embed`** | Perintah *compiler* untuk membungkus file luar (seperti `.sql`) ke dalam program. |
| **`_` (Blank ID)** | Pengganti variabel untuk **mengabaikan** hasil fungsi yang tidak dipakai (wajib di Go). |
| **`*` (Pointer)** | Kunci akses langsung ke alamat memori asli data. |
| **`&` (Address-of)** | Mengambil lokasi memori variabel agar bisa diisi nilainya oleh fungsi (misal: database). |
| **`:=`** | Membuat & mengisi variabel baru secara cepat. |
| **`defer`** | Perintah "tunda dulu" agar kode dijalankan saat fungsi hampir selesai. |
| **CORS** | Keamanan agar website luar boleh akses API kita. |
| **`gin.Context`** | Wadah data permintaan dan jawaban HTTP. |
| **`struct`** | Cetakan objek data. |

---

## 🛠️ 3. Langkah-Langkah Membuat Proyek Go Baru

1. **Inisialisasi**: `mkdir bab6 && cd bab6 && go mod init bab6`
2. **Install Library**:
    * `go get github.com/gin-gonic/gin`: Framework Web Server.
    * `go get github.com/lib/pq`: Driver PostgreSQL untuk Go.
    * `go get github.com/rubenv/sql-migrate`: Alat migrasi database otomatis.
    * `go get github.com/gin-contrib/cors`: Satpam CORS.

---

## 📘 4. Bedah Tuntas Kode (Urutan: Entity ➔ Migration ➔ Database ➔ Repo ➔ Controller ➔ Router ➔ Main)

### 4.1. `entities/card_bridge_entity.go` (Definisi Data)

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

**Penjelasan:**
* `package entities`: File ini masuk dalam kelompok `entities`.
* `type Card struct`: Membuat cetakan data `Card`. Huruf `C` kapital agar *exported* (bisa diakses package lain).
* `ID string`: Field bernama `ID` tipe teks.
* `` `json: "id"` ``: Agar saat dikirim ke format JSON, nama fieldnya menjadi `id` (huruf kecil).

---

### 4.2. `databases/sql_migrations/1_initiate.sql` (Skrip Database)

```sql
-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE card (
    id varchar(20)
);
-- +migrate StatementEnd
```

**Penjelasan:**
* `-- +migrate Up`: Instruksi alat migrasi agar menjalankan perintah ini saat aplikasi mulai.
* `CREATE TABLE ...`: Perintah SQL standar untuk membuat tabel `card`.

---

### 4.3. `databases/database.go` (Koneksi & Migrasi)

```go
package database

import (
	"database/sql"
	"embed"
	"fmt"
	migrate "github.com/rubenv/sql-migrate"
)

//go:embed sql_migrations/*.sql
var dbMigration embed.FS
var DbConnection *sql.DB

func DBMigration(dbParam *sql.DB) {
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigration,
		Root:       "sql_migrations",
	}
	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)
	if errs != nil { panic(errs) }
	DbConnection = dbParam
	fmt.Println("Migrasi berhasil dilakukan", n, migrations)
}
```

**Penjelasan Fungsi `DBMigration(dbParam *sql.DB)`:**
* `dbParam *sql.DB`: Parameter kunci koneksi database yang dikirim dari `main.go`. Tanpa ini, fungsi tidak tahu ke mana harus menjalankan migrasi.
* `migrate.Exec(...)`: Fungsi untuk menjalankan migrasi SQL. `n` (jumlah file sukses), `errs` (pesan error).
* `panic(errs)`: Jika terjadi error fatal saat migrasi (misal database salah), program dihentikan paksa agar developer segera memperbaikinya.
* `DbConnection = dbParam`: Menyimpan kunci koneksi ke variabel global agar bisa dipakai Koki (Repository).

---

### 4.4. `repositories/card_bridge_repo.go` (Query SQL)

```go
package repositories
import ("bab6/entities"; "database/sql")

func GetCards(db *sql.DB) (result []entities.Card, err error) {
	sql := "SELECT * FROM card"
	rows, err := db.Query(sql)
	if err != nil { return }
	defer rows.Close()

	for rows.Next() {
		var data entities.Card
		err = rows.Scan(&data.ID)
		if err != nil { return }
		result = append(result, data)
	}
	return
}

func InsertCard(db *sql.DB, card entities.Card) (err error) {
	sql := "INSERT INTO card(id) values($1)"
	_, err = db.Exec(sql, card.ID)
	return err
}
```

**Penjelasan Fungsi:**
* `GetCards`: Membaca data dengan `Query`. `rows.Scan(&data.ID)`: Data dari database disalin ke `data.ID`. Tanda `&` memberikan lokasi memori variabel `data.ID` agar database tahu di mana harus menaruh datanya.
* `InsertCard`: Menjalankan `db.Exec`.
  * **Simbol `_` (Blank Identifier)**: Fungsi `db.Exec` mengembalikan dua nilai: `sql.Result` dan `error`. Kita hanya butuh `error` untuk cek apakah berhasil atau tidak. Karena kita **tidak butuh** `sql.Result`, kita gunakan **`_`** untuk "membuang" nilai tersebut agar compiler tidak error karena variabel tidak terpakai.

---

### 4.5. `controllers/card_bridge_controllers.go` (HTTP Handler)

```go
package controllers
import ("bab6/databases"; "bab6/entities"; "bab6/repositories"; "net/http"; "github.com/gin-gonic/gin")

func GetCards(c *gin.Context) {
	var result gin.H
	card, err := repositories.GetCards(database.DbConnection)
	if err != nil { result = gin.H{"result": err.Error()} }
	else { result = gin.H{"result": card} }
	c.JSON(http.StatusOK, result)
}
```

**Penjelasan:**
* `c *gin.Context`: Wadah seluruh data Request (permintaan user) dan alat Response (mengirim JSON balik).
* `c.JSON(http.StatusOK, result)`: Mengirim data ke user dalam bentuk JSON. `http.StatusOK` adalah kode HTTP `200` (Sukses).

---

### 4.6. `routers/card_bridge_routers.go` (Routing)

```go
package routers
import ("bab6/controllers"; "time"; "github.com/gin-contrib/cors"; "github.com/gin-gonic/gin")

func StartServer() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{...})) // Middleware CORS
	router.GET("/cards", controllers.GetCards)
	// ... rute lainnya
	return router
}
```

---

### 4.7. `main.go` (Entry Point)

```go
package main
import ("bab6/databases"; "bab6/routers"; "database/sql"; "fmt"; "log"; _ "github.com/lib/pq")

func main() {
	var PORT = ":50000"
	DB, err := sql.Open("postgres", psqlInfo)
	if err != nil { log.Fatal(err) }

	database.DBMigration(DB) // Jalankan migrasi
	defer DB.Close()         // Tutup koneksi saat program mati

	routers.StartServer().Run(PORT) // Jalankan server
}
```

* `_ "github.com/lib/pq"`: Mengimpor driver PostgreSQL secara *side-effect* agar dikenali oleh Go.
* `defer DB.Close()`: Memastikan koneksi database diputus dengan rapi saat aplikasi mati, agar tidak menggantung di memori server.
