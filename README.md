# Dokumentasi Mendalam & Panduan Lengkap REST API Card Bridge (`bab6`)

*Panduan Teknis Sangat Rinci: Bedah Total File, Library, Variabel, Sintaks, dan Fungsi. Dirancang untuk Pemula yang Ingin Memahami "Mengapa" dan "Bagaimana" Kode Golang Bekerja.*

---

## 💡 1. Konsep Dasar & Analogia Sederhana (Untuk Orang Awam)

Aplikasi ini menggunakan **Layered Architecture** (Arsitektur Berlapis). Bayangkan sistem ini seperti sebuah **Restoran Profesional**:

```
[ Pelanggan / Browser / Postman ]
               │
               ▼  1. Pesan/Update Makanan (HTTP Request: GET, POST, PUT, DELETE)
┌──────────────────────────────────────────────┐
│  A. ROUTER (Buku Menu & Kasir)              │ -> Menerima pesanan dan mengarahkan ke Pelayan yang tepat.
└──────────────────────────────────────────────┘
               │
               ▼  2. Meneruskan Pesanan
┌──────────────────────────────────────────────┐
│  B. CONTROLLER (Pelayan / Waiter)           │ -> Memeriksa format pesanan & menyajikan makanan (JSON Response).
└──────────────────────────────────────────────┘
               │
               ▼  3. Minta Bahan / Olah Data
┌──────────────────────────────────────────────┐
│  C. REPOSITORY (Koki / Chef di Dapur)        │ -> Membaca & Menulis data langsung ke Gudang (Database SQL).
└──────────────────────────────────────────────┘
               │
               ▼  4. Ambil / Simpan Data
┌──────────────────────────────────────────────┐
│  D. DATABASE & ENTITY (Gudang & Bentuk Wadah)│ -> PostgreSQL (Gudang Data) & Card Struct (Bentuk Wadah Kartu).
└──────────────────────────────────────────────┘
```

---

## 📖 2. Glosarium Simbol & Istilah Penting Golang

Sebelum membaca kode, pahami simbol ini. Ini adalah "tata bahasa" dalam Go:

| Simbol / Istilah | Penjelasan Sederhana untuk Pemula |
| :--- | :--- |
| **`//go:embed`** | Perintah *compiler* khusus untuk membungkus file eksternal (seperti `.sql`) langsung ke dalam binary aplikasi. |
| **`_` (Blank ID)** | Pengganti variabel untuk **mengabaikan** hasil fungsi yang tidak dipakai. Ini wajib agar tidak *error* saat *compile*. |
| **`*` (Pointer)** | Kunci akses langsung ke alamat memori asli data. Contoh `*sql.DB` artinya kita memegang kunci akses ke koneksi database aslinya. |
| **`&` (Address-of)** | Mengambil lokasi memori variabel agar bisa diisi nilainya secara langsung (misal saat membaca data dari database). |
| **`:=`** | Cara cepat membuat & mengisi variabel baru tanpa `var`. |
| **`[]` (Slice)** | Daftar/List data yang ukurannya bisa fleksibel. |
| **`defer`** | Perintah tunda kode agar dijalankan persis sebelum fungsi selesai. |
| **`import`** | Perintah untuk membawa masuk library agar bisa digunakan di file kita. |
| **CORS** | Keamanan browser agar website luar tidak sembarangan mengambil data API kita. |

---

## 🛠️ 3. Instalasi Proyek (Langkah demi Langkah)

1. **Inisialisasi**:
   ```bash
   mkdir testbab6 && cd testbab6
   go mod init bab6
   ```
2. **Install Library**:
   ```bash
   go get github.com/gin-gonic/gin
   go get github.com/lib/pq
   go get github.com/rubenv/sql-migrate
   ```

---

## 📘 4. Bedah Tuntas File (Sesuai Urutan Alur Sistem)

### 4.1. `entities/card_bridge_entity.go` (Definisi Data)

**Tugas:** Wadah data (Struct) untuk objek kartu.

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### Bedah File:
* **Library Import**: Tidak ada.
* **Penjelasan Kode**:
  * `package entities`: Menandakan file ini milik folder `entities`.
  * `type Card struct`: Membuat cetakan data bernama `Card`.
  * `ID string`: Wadah teks bernama `ID`.
  * `` `json: "id"` ``: Agar saat dikirim ke format JSON, nama fieldnya menjadi `id`.

---

### 4.2. `databases/sql_migrations/1_initiate.sql` (Skrip Database)

**Tugas:** Perintah SQL pembuat tabel.

```sql
-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE card (id varchar(20));
-- +migrate StatementEnd
```

#### Bedah File:
* **Library Import**: Tidak ada.
* **Penjelasan**:
  * `-- +migrate Up`: Instruksi alat migrasi agar menjalankan perintah ini saat aplikasi mulai.
  * `CREATE TABLE card (id varchar(20));`: SQL murni untuk membuat tabel `card`.

---

### 4.3. `databases/database.go` (Koneksi & Migrasi)

**Tugas:** Koneksi database & migrasi otomatis.

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

#### Bedah Fungsi `DBMigration(dbParam *sql.DB)`:
* **`dbParam *sql.DB`**: Kunci koneksi database yang dikirim dari `main.go`.
* **`dbMigration`**: Variabel sistem file virtual hasil *embed* folder migrasi.
* **`migrate.Exec(...)`**: Fungsi inti untuk menjalankan skrip SQL.
  * `n`: Jumlah file SQL yang berhasil dijalankan.
  * `errs`: Pesan error jika gagal.
* **`panic(errs)`**: Jika ada error, hentikan program paksa.

---

### 4.4. `repositories/card_bridge_repo.go` (Query SQL)

**Tugas:** Eksekusi perintah SQL (CRUD).

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
	_, err = db.Exec(sql, card.ID) // _ mengabaikan metadata
	return err
}
```

#### Bedah Fungsi:
* `GetCards`: Membaca data dengan `Query`. `rows.Scan(&data.ID)` memindahkan data dari database ke variabel Go.
* `InsertCard`: Menggunakan `db.Exec`. Tanda **`_`** digunakan untuk **mengabaikan** metadata hasil eksekusi (karena kita hanya butuh error-nya). Jika tidak pakai `_`, Go akan *error* karena variabel metadata tidak dipakai.

---

### 4.5. `controllers/card_bridge_controllers.go` (HTTP Handler)

**Tugas:** Menangani permintaan HTTP (masuk & keluar).

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

#### Bedah Fungsi:
* `GetCards`: Memanggil `repositories.GetCards` dengan kunci `database.DbConnection`.
* `c.JSON(...)`: Fungsi Gin untuk mengirim balasan dalam format JSON.

---

### 4.6. `routers/card_bridge_routers.go` (Alamat URL)

**Tugas:** Menentukan jalur URL.

```go
package routers
import ("bab6/controllers"; "time"; "github.com/gin-contrib/cors"; "github.com/gin-gonic/gin")

func StartServer() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{...}))
	router.GET("/cards", controllers.GetCards)
	// ... rute lainnya
	return router
}
```
* **Penjelasan**: `cors.New(...)` adalah "Satpam" agar domain luar bisa akses API ini.

---

### 4.7. `main.go` (Entry Point)

**Tugas:** Titik masuk utama aplikasi.

```go
package main
import (
	database "bab6/databases"; "bab6/routers"; "database/sql"; "fmt"; "log"; _ "github.com/lib/pq"
)

func main() {
	var PORT = ":50000"
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil { log.Fatal(err) }

	database.DBMigration(DB) // Jalankan migrasi
	defer DB.Close()         // Tutup koneksi saat program mati

	routers.StartServer().Run(PORT) // Jalankan server
}
```

---

## 🚀 Cara Menjalankan
1. `CREATE DATABASE mcsbab6;`
2. `go mod tidy`
3. `go run main.go`
