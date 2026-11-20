# 📘 Dokumentasi Utama & Bedah Teknis REST API Card Bridge (`bab6`)

Dokumentasi ini dirancang sebagai panduan pembelajaran komprehensif untuk memahami cara membangun REST API menggunakan bahasa pemrograman Go (Golang). Kita akan membedah proyek ini dari nol, menjelaskan filosofi di balik setiap baris kode, hingga memahami simbol-simbol teknis yang sering membingungkan pemula.

---

## 💡 1. Penjelasan Proyek (Overview)

Aplikasi **Card Bridge (bab6)** adalah sebuah sistem backend sederhana namun kokoh yang berfungsi untuk mengelola identitas kartu. Proyek ini menggunakan **Arsitektur Berlapis (Layered Architecture)**, sebuah standar industri yang memisahkan kode berdasarkan tanggung jawabnya agar aplikasi mudah dikelola, diuji, dan dikembangkan di masa depan.

**Tujuan Proyek:**
*   Menyediakan endpoint API untuk melihat daftar kartu (`GET`).
*   Menyediakan endpoint API untuk menambah kartu baru (`POST`).
*   Menyediakan endpoint API untuk menghapus data kartu (`DELETE`).
*   Mengotomatiskan pembuatan tabel database menggunakan sistem migrasi.

---

## 🍴 2. Analogi Sistem: Restoran Digital

Agar lebih mudah membayangkan alur data, mari gunakan analogi sebuah restoran:

*   **PELANGGAN (Client/Postman)**: Pihak yang lapar akan data. Mereka mengirimkan pesanan melalui HTTP Request (GET, POST, dll).
*   **KASIR (Router)**: Orang pertama yang menerima pesanan di depan. Kasir mengecek menu (URL) dan mengarahkan pesanan ke pelayan yang tepat.
*   **PELAYAN (Controller)**: Menemui pelanggan, mencatat detail pesanan, memastikan formatnya benar, dan berkomunikasi dengan bagian dapur. Setelah makanan siap, pelayan menyajikannya dalam piring JSON.
*   **KOKI (Repository)**: Sang ahli dapur. Koki tidak tahu siapa pelanggan di depan; dia hanya fokus mengolah bahan mentah menjadi masakan berdasarkan perintah (Query SQL).
*   **GUDANG BAHAN (Database)**: Tempat penyimpanan permanen seluruh bahan baku (Data di PostgreSQL).
*   **CETAKAN PIRING (Entity)**: Standar bentuk makanan yang akan disajikan. Ini memastikan semua data kartu yang keluar dari sistem memiliki bentuk yang konsisten.

---

## 📖 3. Glosarium: Simbol, Sintaks, & Istilah Teknis

Berikut adalah kamus kecil untuk memahami simbol-simbol khusus di dalam kode Go:

| Simbol / Istilah | Nama Teknis | Kegunaan & Penjelasan |
| :--- | :--- | :--- |
| **`//go:embed`** | **Compiler Directive** | Memberitahu program untuk "mencaplok" file luar (seperti file SQL) dan memasukkannya ke dalam file aplikasi akhir sehingga file tersebut tidak akan hilang saat aplikasi dipindah-pindah. |
| **`_`** | **Blank Identifier** | "Variabel Tempat Sampah". Digunakan untuk menampung hasil fungsi yang **tidak ingin kita pakai**. Karena Go melarang adanya variabel menganggur, kita gunakan `_` untuk membuangnya. |
| **`*`** | **Pointer** | Menunjuk langsung ke "alamat asli" data di memori komputer, bukan sekadar fotokopinya. Menggunakan pointer menghemat memori karena kita tidak menyalin data besar berulang-ulang. |
| **`&`** | **Address-of Operator** | Digunakan untuk mendapatkan "nomor rumah" atau alamat memori dari sebuah variabel. |
| **`:=`** | **Short Declaration** | Cara singkat membuat variabel baru sekaligus memberi isi tanpa perlu menulis kata `var`. |
| **`defer`** | **Defer Statement** | Perintah untuk menunda aksi sampai fungsi benar-benar berakhir. Biasanya dipakai untuk "menutup pintu" (koneksi database) agar tidak terjadi kebocoran memori. |
| **`struct`** | **Structure** | Wadah untuk mengelompokkan berbagai data menjadi satu objek. |
| **`package`** | **Paket** | Cara Go mengelompokkan file-file yang memiliki fungsi serupa dalam satu folder. |
| **CORS** | **Cross-Origin Resource Sharing** | Protokol keamanan yang mengizinkan website dari luar (seperti aplikasi Frontend) untuk mengakses API kita. |
| **`gin.Context`** | **Konteks Gin** | Objek "sakti" yang membawa semua informasi tentang permintaan user (Header, Body, URL) dan alat untuk mengirim balik jawaban. |

---

## 🛠️ 4. Langkah Membuat Proyek & Instalasi Library

Untuk membangun aplikasi ini dari awal, ikuti langkah-langkah berikut:

### A. Inisialisasi Proyek
```bash
# 1. Buat folder proyek
mkdir testbab6
cd testbab6

# 2. Inisialisasi modul Go (Nama modul 'bab6' akan dipakai saat import)
go mod init bab6
```

### B. Instalasi Library Eksternal
Library eksternal adalah alat bantu yang dibuat oleh komunitas untuk mempermudah pekerjaan kita:

1.  **Framework Gin** (`go get github.com/gin-gonic/gin`)
    *   **Kegunaan**: Mesin penggerak web server. Menangani rute URL, request, dan response secara efisien.
2.  **Driver PostgreSQL** (`go get github.com/lib/pq`)
    *   **Kegunaan**: Jembatan komunikasi antara bahasa Go dengan database PostgreSQL.
3.  **Sql-Migrate** (`go get github.com/rubenv/sql-migrate`)
    *   **Kegunaan**: Alat untuk manajemen struktur tabel database agar sinkron antara kode dan database asli.
4.  **CORS Middleware** (`go get github.com/gin-contrib/cors`)
    *   **Kegunaan**: Mengatur izin akses keamanan API.

---

## 📘 5. Bedah Tuntas Kode (File demi File)

Berikut adalah bedah total setiap file kodenya secara berurutan:

---

### 5.1. `entities/card_bridge_entity.go` (Cetak Biru Data)
**Tugas**: Mendefinisikan struktur data kartu yang akan diolah di seluruh aplikasi.

#### KODE:
```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### PENJELASAN MENDETAIL:
*   **`package entities`**: Baris ini memberi tahu Go bahwa file ini adalah bagian dari folder `entities`.
*   **`type Card struct`**: Kita menciptakan tipe data baru bernama `Card`. Nama ini menggunakan huruf **Kapital (C)** agar bisa di-import oleh folder lain (Public).
*   **`ID string`**: Di dalam kartu, terdapat data `ID` yang bertipe teks (`string`).
*   **`` `json: "id"` ``**: Ini disebut **Struct Tag**. Fungsinya: Saat data dikirim ke pelanggan dalam format JSON, nama field yang tadinya `ID` (Besar) akan otomatis berubah menjadi `id` (kecil) sesuai standar web.

---

### 5.2. `databases/sql_migrations/1_initiate.sql` (Instruksi Database)
**Tugas**: Skrip SQL yang akan dijalankan oleh sistem migrasi untuk membuat tabel.

#### KODE:
```sql
-- +migrate Up
-- +migrate StatementBegin

CREATE TABLE card (
    id varchar(20)
);

-- +migrate StatementEnd
```

#### PENJELASAN MENDETAIL:
*   **`-- +migrate Up`**: Penanda bahwa skrip ini harus dijalankan saat sistem melakukan *update* ke atas (membuat tabel).
*   **`-- +migrate StatementBegin/End`**: Membungkus perintah SQL agar dijalankan sebagai satu kesatuan perintah yang utuh.
*   **`CREATE TABLE card`**: Perintah standar SQL untuk membuat tabel baru bernama `card`.
*   **`id varchar(20)`**: Membuat satu kolom bernama `id` dengan tipe teks dinamis maksimal 20 karakter.

---

### 5.3. `databases/database.go` (Pusat Koneksi)
**Tugas**: Menangani proses pembukaan gerbang database dan menjalankan migrasi otomatis saat aplikasi pertama kali menyala.

#### KODE:
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

	if errs != nil {
		panic(errs)
	}
	DbConnection = dbParam
	fmt.Println("Migrasi berhasil dilakukan", n, migrations)
}
```

#### PENJELASAN MENDETAIL:
*   **`//go:embed sql_migrations/*.sql`**: Mengambil semua file SQL di folder migrasi dan memasukkannya ke dalam variabel `dbMigration` bertipe `embed.FS`.
*   **`var DbConnection *sql.DB`**: Variabel global untuk menampung koneksi database. Kita beri nama `DbConnection` dengan huruf **D** besar agar bisa dipanggil dari folder Controller dan Repository.
*   **`func DBMigration(dbParam *sql.DB)`**: Fungsi ini menerima bahan baku berupa koneksi database (`dbParam`).
    *   **`migrations := ...`**: Menyiapkan sumber file SQL dari folder `sql_migrations`.
    *   **`migrate.Exec(...)`**: Perintah untuk mengeksekusi file SQL ke database PostgreSQL.
    *   **`n`**: Berisi angka jumlah file yang sukses dijalankan.
    *   **`errs`**: Menampung pesan error jika proses gagal.
    *   **`panic(errs)`**: Jika terjadi error (misal database salah password), hentikan program seketika dengan pesan error agar developer tahu ada masalah besar.

---

### 5.4. `repositories/card_bridge_repo.go` (Dapur Rekaman)
**Tugas**: Tempat Koki bekerja. Semua perintah SQL (`SELECT`, `INSERT`) ditulis di sini.

#### KODE:
```go
package repositories

import (
	"bab6/entities"
	"database/sql"
)

func GetCards(db *sql.DB) (result []entities.Card, err error) {
	sql := "SELECT * FROM card"
	rows, err := db.Query(sql)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data entities.Card
		err = rows.Scan(&data.ID)
		if err != nil {
			return
		}
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

#### PENJELASAN MENDETAIL FUNGSI:
1.  **Fungsi `GetCards`**: Mengambil semua data kartu.
    *   **`db.Query(sql)`**: Menjalankan perintah cari. Hasilnya disimpan di variabel `rows`.
    *   **`defer rows.Close()`**: Memastikan data yang sedang dibaca ditutup setelah selesai, agar koneksi tidak "nyangkut".
    *   **`for rows.Next()`**: Melakukan perulangan untuk setiap baris data yang ditemukan di database.
    *   **`rows.Scan(&data.ID)`**: Memindahkan data dari database ke variabel `data.ID`. Simbol **`&`** digunakan agar database tahu lokasi memori pasti di mana dia harus meletakkan datanya.
    *   **`append(result, data)`**: Menambahkan satu kartu ke dalam daftar panjang `result`.
2.  **Fungsi `InsertCard`**: Menambah kartu baru.
    *   **`db.Exec(sql, card.ID)`**: Menjalankan perintah tambah.
    *   **`_` (Blank Identifier)**: Karena `db.Exec` memberikan info statistik (berapa baris yang berubah) yang tidak kita butuhkan, kita gunakan `_` untuk membuangnya.
  3.  **Fungsi `DeleteCard`**: Menghapus kartu berdasarkan ID.
      *   **`c.Param("id")`**: Mengambil nilai ID dari URL (misal: `/card/delete/ABC-123`).
      *   **`repositories.DeleteCard(...)`**: Memanggil repository untuk menjalankan perintah `DELETE FROM card WHERE id = $1`.
      *   **`c.JSON(http.StatusOK, ...)`**: Mengirim respons sukses ke user dengan pesan "Data Berhasil Dihapus" dan ID yang dihapus.

---

### 5.5. `controllers/card_bridge_controllers.go` (Pelayan Restoran)
**Tugas**: Menerima request dari user, memanggil Repository, dan memberikan jawaban (Response).

#### KODE:
```go
package controllers

import (
	database "bab6/databases"
	"bab6/entities"
	"bab6/repositories"
	"net/http"
	"github.com/gin-gonic/gin"
)

func GetCards(c *gin.Context) {
	var result gin.H
	card, err := repositories.GetCards(database.DbConnection)

	if err != nil {
		result = gin.H{"result": err.Error()}
	} else {
		result = gin.H{"result": card}
	}
	c.JSON(http.StatusOK, result)
}

func InsertCard(c *gin.Context) {
	var card entities.Card
	idCard := c.Param("id")
	card.ID = idCard

	err := repositories.InsertCard(database.DbConnection, card)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, card)
}
```

#### PENJELASAN MENDETAIL:
*   **`c *gin.Context`**: Parameter wajib di Gin. Menampung data masuk dan alat kirim data keluar.
*   **`idCard := c.Param("id")`**: Mengambil teks ID yang diketik user di alamat URL (misal: `/card/input/ABC-123`, maka ID-nya adalah `ABC-123`).
*   **`c.JSON(...)`**: Mengirimkan jawaban balik ke user.
    *   `http.StatusOK`: Kode angka `200` (Berhasil).
    *   `http.StatusInternalServerError`: Kode angka `500` (Ada masalah di server/database).

---

### 5.6. `routers/card_bridge_routers.go` (Daftar Menu URL)
**Tugas**: Mengatur rute jalan raya. Alamat URL mana akan diproses oleh fungsi mana.

#### KODE:
```go
package routers

import (
	"bab6/controllers"
	"time"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		MaxAge:           12 * time.Hour,
	}))

	router.GET("/cards", controllers.GetCards)
	router.POST("/card/input/:id", controllers.InsertCard)
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	return router
}
```

#### PENJELASAN MENDETAIL:
*   **`router.Use(cors.New(...))`**: Memasang aturan keamanan CORS agar aplikasi frontend dari website manapun (`"*"`) bisa mengambil data dari API ini.
*   **`router.GET("/cards", ...)`**: Menentukan bahwa alamat `GET /cards` akan dilayani oleh fungsi `GetCards` di folder controller.
*   **`:id`**: Ini adalah variabel dinamis. Artinya bagian ini bisa diisi apa saja oleh user di URL.

---

### 5.7. `main.go` (Jantung Aplikasi)
**Tugas**: Tempat tombol ON/OFF. Menyiapkan semua konfigurasi dan menjalankan server.

#### KODE:
```go
package main

import (
	database "bab6/databases"
	"bab6/routers"
	"database/sql"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "password_anda"
	dbName   = "mcsbab6"
)

func main() {
	var PORT = ":50000"

	psqlInfo := fmt.Sprintf(
		`host=%s port=%d user=%s password=%s dbname=%s sslmode=disable`,
		host, port, user, password, dbName,
	)

	DB, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Error membuka DB: ", err)
	}

	database.DBMigration(DB)
	defer DB.Close()

	fmt.Printf("Aplikasi berjalan di port %s\n", PORT)
	routers.StartServer().Run(PORT)
}
```

#### PENJELASAN MENDETAIL:
*   **`_ "github.com/lib/pq"`**: Mengimpor driver Postgres. Tanda `_` (Side-effect import) berarti kita hanya ingin menjalankan kode inisialisasi di dalam library tersebut tanpa memanggil fungsinya secara langsung.
*   **`psqlInfo := fmt.Sprintf(...)`**: Merangkai semua variabel (host, user, pass) menjadi satu kalimat perintah koneksi yang dimengerti PostgreSQL.
*   **`sql.Open("postgres", ...)`**: Mengetuk pintu database untuk memulai koneksi.
*   **`log.Fatal(...)`**: Jika terjadi error fatal di awal, hentikan program dan catat pesannya.
*   **`routers.StartServer().Run(PORT)`**: Baris perintah paling sakti. Menyalakan mesin server web dan mulai menunggu pesanan user di port `50000`.

---

## 🚀 6. Cara Menjalankan

1.  **Database**: Pastikan PostgreSQL menyala dan buat database kosong bernama `mcsbab6`.
2.  **Mod Tidy**: Jalankan `go mod tidy` di terminal untuk merapikan semua library.
3.  **Jalankan**: Ketik `go run main.go`.
4.  **Tes Semua Endpoint di Postman**:
    *   **GET** `http://localhost:50000/cards` → Ambil daftar kartu.
    *   **POST** `http://localhost:50000/card/input/{id}` → Tambah kartu baru (ganti `{id}` dengan nilai contoh).
    *   **DELETE** `http://localhost:50000/card/delete/{id}` → Hapus kartu berdasarkan ID (ganti `{id}` dengan nilai yang ada).

---
*Dokumentasi ini dibuat untuk memastikan Anda memahami filosofi dan teknis dari setiap baris kode yang Anda miliki.*
