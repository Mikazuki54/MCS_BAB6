# 📘 Dokumentasi Super Lengkap & Panduan Utama REST API Card Bridge (`bab6`)

Selamat datang di panduan teknis terdahulu! Dokumentasi ini dirancang bukan sekadar untuk dibaca, tapi untuk dipelajari baris demi baris. Kami akan membedah setiap file, setiap fungsi, setiap variabel, bahkan setiap simbol unik yang ada di dalam proyek Go ini. 

---

## 💡 1. Konsep Dasar & Arsitektur (Analogi Restoran)

Untuk memudahkan pemula, bayangkan aplikasi ini adalah sebuah **Restoran Bintang Lima**. Dalam Go, kita membagi tugas agar kode tidak berantakan (**Layered Architecture**):

1.  **PELANGGAN (Client/Postman)**: Orang yang memesan makanan (mengirim HTTP Request).
2.  **KASIR / MENU (Router)**: Tempat pelanggan memilih apa yang ingin dipesan (URL) dan mengarahkan ke pelayan yang tepat.
3.  **PELAYAN / WAITER (Controller)**: Menanyakan detail pesanan, memastikan formatnya benar, dan menyerahkan piring hasil masakan (JSON Response) kembali ke pelanggan.
4.  **KOKI / CHEF (Repository)**: Orang yang tahu resep (Query SQL) dan mengambil bahan masakan dari gudang.
5.  **GUDANG BAHAN (Database)**: Tempat penyimpanan semua data (PostgreSQL).
6.  **CETAKAN PIRING (Entity)**: Bentuk atau wadah makanan yang disajikan agar seragam.

---

## 📖 2. Glosarium Simbol, Sintaks, & Istilah Teknis

Sebelum masuk ke kode, pastikan Anda memahami "kamus" simbol Golang berikut:

| Simbol / Sintaks | Nama | Penjelasan Spesifik untuk Pemula |
| :--- | :--- | :--- |
| **`package`** | Paket | Kelompok folder tempat file berada. Semua file di folder yang sama harus punya nama package yang sama. |
| **`import`** | Impor | Perintah untuk memanggil "alat bantu" (library) dari folder lain atau dari internet. |
| **`func`** | Fungsi | Blok kode yang menjalankan tugas tertentu. |
| **`_`** | **Blank Identifier** | Variabel "tempat sampah". Jika sebuah fungsi menghasilkan nilai tapi kita tidak mau menggunakannya, kita pakai `_` agar Go tidak komplain "variabel tidak terpakai". |
| **`*`** | **Pointer** | Alamat rumah data. `*sql.DB` berarti kita memegang alamat asli koneksi database, bukan fotokopinya. |
| **`&`** | **Address-of** | Mengambil kunci alamat. `&data.ID` berarti kita memberikan lokasi memori variabel ID agar bisa diisi oleh database. |
| **`:=`** | **Short Declaration** | Membuat variabel baru + mengisi nilainya secara otomatis tanpa mengetik `var`. |
| **`[]`** | **Slice** | Sebuah daftar (list) data yang panjangnya bisa berubah-ubah. |
| **`defer`** | **Penundaan** | Perintah "nanti ya dijalankannya". Biasanya dipakai untuk memastikan pintu koneksi ditutup tepat sebelum fungsi selesai. |
| **`struct`** | **Struktur** | Kumpulan beberapa variabel yang dibungkus menjadi satu objek (misal: Kartu punya ID). |

---

## 🛠️ 3. Langkah-Langkah Membuat Proyek Go & Gin Baru (Setup Dari Nol)

Jika Anda ingin membuat proyek seperti ini dari awal, ikuti langkah-langkah ini:

### Langkah 1: Inisialisasi Folder
Buka terminal, buat folder, dan jalankan perintah ini untuk membuat file identitas proyek (`go.mod`):
```bash
mkdir nama-proyek
cd nama-proyek
go mod init nama-modul
```

### Langkah 2: Instalasi Library Utama
Jalankan perintah ini satu per satu. Berikut penjelasan fungsinya:

1.  **Framework Gin (`github.com/gin-gonic/gin`)**
    *   **Fungsi**: Mesin utama untuk membuat Web Server. Gin memudahkan kita mengatur alamat URL, membaca data kiriman user, dan mengirim jawaban JSON.
2.  **Driver PostgreSQL (`github.com/lib/pq`)**
    *   **Fungsi**: "Sopir" atau penghubung agar bahasa Go bisa mengobrol dengan database PostgreSQL. Tanpa ini, Go tidak akan mengerti perintah database.
3.  **Library Migrasi (`github.com/rubenv/sql-migrate`)**
    *   **Fungsi**: Alat otomatis untuk membuat tabel di database hanya dengan file teks `.sql`.
4.  **Middleware CORS (`github.com/gin-contrib/cors`)**
    *   **Fungsi**: Satpam keamanan. Mengizinkan website frontend (seperti React/Vue) untuk mengambil data dari API kita tanpa diblokir oleh browser.

---

## 📘 4. Bedah Tuntas Kode Per File (Baris demi Baris)

Kami akan mengurutkan penjelasan sesuai alur jalannya data:

### 4.1. `entities/card_bridge_entity.go` (Definisi Wadah Data)

File ini adalah cetakan dasar. Ibarat piranti makan, ini menentukan piring mana yang dipakai untuk menyajikan kartu.

```go
package entities // File ini masuk ke kelompok entities

// Card adalah tipe data bentukan (struct) yang mewakili 1 kartu
type Card struct {
	ID string `json: "id"` // Properti ID bertipe teks (string)
}
```

*   **`type Card struct`**: Membuat objek baru bernama `Card`. Huruf `C` besar artinya objek ini "Exported" (bisa dipakai oleh folder lain).
*   **`ID string`**: Variabel di dalam objek bernama `ID`.
*   **`` `json: "id"` ``**: Ini disebut **Struct Tag**. Fungsinya memberitahu Go: "Kalau data ini dikirim ke user lewat JSON, ganti tulisan `ID` jadi `id` huruf kecil ya".

---

### 4.2. `databases/sql_migrations/1_initiate.sql` (Resep Struktur Database)

File ini bukan kode Go, tapi kode SQL untuk membangun gudang data kita.

```sql
-- +migrate Up
-- Perintah untuk sql-migrate: "Jalankan ini kalau aplikasi baru nyala"

-- +migrate StatementBegin
-- Mulai satu paket perintah agar tidak terputus di tengah jalan

CREATE TABLE card (
    id varchar(20) -- Buat tabel bernama 'card' dengan kolom 'id' teks max 20 karakter
);

-- +migrate StatementEnd
-- Selesai paket perintah
```

---

### 4.3. `databases/database.go` (Pengelola Migrasi & Koneksi Global)

File ini bertugas menjalankan resep SQL di atas secara otomatis.

```go
package database

import (
	"database/sql" // Pustaka standar Go untuk database
	"embed"        // Fitur untuk membungkus file luar ke dalam aplikasi
	"fmt"          // Pustaka untuk mencetak tulisan ke terminal
	migrate "github.com/rubenv/sql-migrate" // Pustaka migrasi pihak ketiga
)

//go:embed sql_migrations/*.sql
// Ini adalah Compiler Directive (Perintah khusus). 
// Gunanya membungkus semua file .sql di folder migrations agar masuk ke variabel dbMigration.
var dbMigration embed.FS

// DbConnection adalah variabel global yang menyimpan kunci akses database. 
// Tanda * artinya pointer (alamat memori asli).
var DbConnection *sql.DB

// Fungsi DBMigration dijalankan untuk mengeksekusi file SQL tadi.
func DBMigration(dbParam *sql.DB) {
	// 1. Menentukan di mana file SQL berada (dalam variabel dbMigration)
	migrations := &migrate.EmbedFileSystemMigrationSource{
		FileSystem: dbMigration,
		Root:       "sql_migrations",
	}

	// 2. Mengeksekusi migrasi ke PostgreSQL
	// n = menampung angka berapa file yang sukses dijalankan.
	// errs = menampung pesan error kalau gagal.
	// dbParam = Kunci koneksi yang dikirim dari main.go.
	// migrate.Up = Perintah untuk "Membuat/Menambah" tabel.
	n, errs := migrate.Exec(dbParam, "postgres", migrations, migrate.Up)

	// 3. Cek Error
	if errs != nil {
		panic(errs) // Kalau error, matikan aplikasi dan kasih tahu penyebabnya
	}

	// 4. Menyimpan koneksi ke variabel global agar bisa dipakai Koki (Repository)
	DbConnection = dbParam

	fmt.Println("Migrasi berhasil dilakukan", n, migrations)
}
```

---

### 4.4. `repositories/card_bridge_repo.go` (Koki - Pengolah Data Langsung)

Di sinilah perintah SQL yang rumit dijalankan.

```go
package repositories

import (
	"bab6/entities" // Mengambil cetakan Card dari folder entities
	"database/sql"  // Mengambil alat bantu SQL
)

// Fungsi GetCards mengambil semua data kartu
func GetCards(db *sql.DB) (result []entities.Card, err error) {
	sql := "SELECT * FROM card" // Teks perintah SQL
	
	// db.Query mengirim perintah SELECT dan mendapatkan banyak baris data (rows)
	rows, err := db.Query(sql)

	if err != nil {
		return // Jika gagal, langsung kembalikan pesan error
	}

	// defer rows.Close() memastikan "pintu gudang data" ditutup saat fungsi selesai.
	// Ini sangat penting agar RAM komputer tidak penuh.
	defer rows.Close()

	// Looping (perulangan) untuk membaca data baris demi baris
	for rows.Next() {
		var data entities.Card // Buat wadah sementara
		
		// rows.Scan memindahkan isi data dari kolom database ke variabel Go.
		// Tanda &data.ID artinya "isi ke alamat memori variabel ID".
		err = rows.Scan(&data.ID)
		
		if err != nil {
			return
		}
		
		// append memasukkan data yang baru dibaca ke dalam daftar list (result)
		result = append(result, data)
	}
	return // Kirim hasil daftar kartu
}

// Fungsi InsertCard untuk menambah kartu baru
func InsertCard(db *sql.DB, card entities.Card) (err error) {
	sql := "INSERT INTO card(id) values($1)" // $1 adalah pengaman (placeholder)
	
	// db.Exec dipakai untuk perintah yang MENGUBAH data (Insert, Update, Delete)
	// Tanda _ (Blank Identifier) dipakai karena db.Exec mengembalikan 'Result' 
	// yang tidak kita butuhkan, jadi kita buang ke "tempat sampah".
	_, err = db.Exec(sql, card.ID)
	
	return err
}

// UpdateCard & DeleteCard fungsinya serupa dengan InsertCard (menggunakan db.Exec)
```

---

### 4.5. `controllers/card_bridge_controllers.go` (Pelayan - Pengatur Permintaan)

Tugasnya menerima input dari user dan memanggil Koki.

```go
package controllers

import (
	database "bab6/databases" // Mengambil variabel global DbConnection
	"bab6/entities"           // Mengambil struct Card
	"bab6/repositories"       // Mengambil fungsi Koki (GetCards, dll)
	"net/http"                // Pustaka untuk kode status HTTP (200, 404, 500)
	"github.com/gin-gonic/gin" // Pustaka Web Framework
)

func GetCards(c *gin.Context) {
	// 1. Panggil Koki (Repository) untuk mengambil data
	card, err := repositories.GetCards(database.DbConnection)

	if err != nil {
		// Jika koki gagal, kirim status 500 (Server Error)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	// 2. Jika berhasil, sajikan piring JSON status 200 (OK)
	// gin.H adalah cara cepat membuat objek JSON
	c.JSON(http.StatusOK, gin.H{"result": card})
}

func InsertCard(c *gin.Context) {
	var card entities.Card
	
	// c.Param("id") mengambil nilai ID dari alamat URL (misal: /input/CARD-01)
	idCard := c.Param("id")
	card.ID = idCard

	// Panggil koki untuk simpan ke database
	err := repositories.InsertCard(database.DbConnection, card)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}
```

---

### 4.6. `routers/card_bridge_routers.go` (Kasir - Pengatur Rute & Keamanan)

Di sini kita menentukan alamat URL dan memasang Satpam (CORS).

```go
package routers

import (
	"bab6/controllers"
	"time"
	"github.com/gin-contrib/cors" // Library Satpam CORS
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	// 1. Inisialisasi Mesin Gin
	router := gin.Default()

	// 2. Memasang Aturan CORS (Satpam)
	// Ini penting agar aplikasi web lain (React/Vue) bisa mengakses API kita.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Mengizinkan semua domain
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour, // Izin berlaku selama 12 jam
	}))

	// 3. Daftar Menu Alamat (Endpoint)
	router.GET("/cards", controllers.GetCards)               // Alamat Lihat Data
	router.POST("/card/input/:id", controllers.InsertCard)   // Alamat Tambah Data
	router.DELETE("/card/delete/:id", controllers.DeleteCard) // Alamat Hapus Data

	return router // Kembalikan mesin yang sudah siap
}
```

---

### 4.7. `main.go` (Pintu Utama - Menyalakan Semua Sistem)

File ini adalah tombol **ON** untuk seluruh aplikasi.

```go
package main

import (
	database "bab6/databases"
	"bab6/routers"
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq" // Side-effect import: Mendaftarkan "sopir" PostgreSQL
)

// Data konfigurasi database (Ganti sesuai DB Anda)
const (
	password = "nickytirta54" 
	dbName   = "mcsbab6"      
)

func main() {
	var PORT = ":50000" // Alamat port pintu masuk server

	// Rangkaian kalimat untuk menyambung ke PostgreSQL
	psqlInfo := fmt.Sprintf("host=localhost port=5432 user=postgres password=%s dbname=%s sslmode=disable", password, dbName)

	// 1. Membuka gerbang koneksi ke database
	// DB = variabel penyimpan koneksi. err = variabel penyimpan pesan error.
	DB, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Gagal membuka DB: ", err) // Jika gagal parah, matikan program
	}

	// 2. Jalankan Migrasi Otomatis (Membuat Tabel)
	// Kita mengirim variabel 'DB' sebagai parameter 'dbParam' ke fungsi DBMigration.
	database.DBMigration(DB)

	// 3. Defer menutup koneksi DB saat server dimatikan
	defer DB.Close()

	// 4. Nyalakan Server Web!
	fmt.Printf("Server berjalan di port %s\n", PORT)
	routers.StartServer().Run(PORT) 
}
```

---

## 🚀 5. Cara Menjalankan Aplikasi

1.  Buka **PostgreSQL** dan buat database kosong:
    `CREATE DATABASE mcsbab6;`
2.  Buka terminal di folder proyek ini.
3.  Jalankan perintah untuk merapikan pustaka:
    `go mod tidy`
4.  Nyalakan aplikasi:
    `go run main.go`
5.  Selesai! Anda bisa mencoba mengakses `http://localhost:50000/cards` di browser.

---

## 🧪 6. Cara Mencoba API (Testing)

Anda bisa menggunakan aplikasi **Postman** atau perintah **cURL** di terminal:

*   **Tambah Kartu**: `POST http://localhost:50000/card/input/CARD-01`
*   **Lihat Semua**: `GET http://localhost:50000/cards`
*   **Hapus Kartu**: `DELETE http://localhost:50000/card/delete/CARD-01`
