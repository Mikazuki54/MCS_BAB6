# Panduan & Dokumentasi Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Spesifik per Potongan Kode, Baris demi Baris, Variabel, Sintaks, dan Fungsi.*

---

## 💡 1. Konsep Dasar & Analogia Sederhana (Untuk Orang Awam)

Sebelum masuk ke kode, mari kita pahami apa yang dibuat oleh aplikasi ini dengan analogi sederhana: **Restoran**.

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

* **REST API**: Cara aplikasi saling berkomunikasi menggunakan protokol HTTP. Klien mengirim *Request* (permintaan) dan server memberikan *Response* (jawaban berupa teks JSON).
* **HTTP Methods**:
  * **`GET`**: Membaca / mengambil data (seperti melihat menu).
  * **`POST`**: Menambahkan data baru (seperti memesan makanan baru).
  * **`PUT`**: Memperbarui/mengubah data yang *sudah ada* (seperti mengganti isi pesanan yang salah).
  * **`DELETE`**: Menghapus data (seperti membatalkan pesanan).
* **Layered Architecture**: Memisah-misah kode ke beberapa folder agar jika ada perubahan atau kerusakan, kita tahu persis folder mana yang harus diperbaiki tanpa mengganggu bagian lain.

---

## 📖 2. Glosarium Istilah Pemula

*Baca ini agar lebih mudah memahami istilah teknis di bawah:*

| Istilah | Penjelasan Sederhana untuk Pemula |
| :--- | :--- |
| **Golang / Go** | Bahasa pemrograman yang cepat, ringan, dan cocok untuk backend server. |
| **Gin Gonic** | Alat bantu (framework) untuk membuat web server di Go. |
| **Struct** | Fitur Go untuk membuat wadah tipe data sendiri (seperti objek). |
| **Pointer (`*`)** | Tanda penunjuk alamat memori komputer asli, bukan menduplikat datanya. |
| **Parameter** | Nilai/bahan masukan yang dikirim ke dalam sebuah fungsi agar bisa diolah. |
| **JSON** | Format teks pertukaran data di internet (`"key": "value"`). |
| **CORS** | Sistem keamanan browser agar web lain tidak sembarangan mengambil data API kita. |
| **`defer`** | Perintah "nanti saja" untuk menjalankan fungsi tepat sebelum fungsi utama selesai. |

---

## 🛠️ 3. Langkah-Langkah Membuat Proyek Go Baru Menggunakan Framework Gin & Library Pendukung

### Langkah 1: Install Golang
Pastikan Go sudah terinstal di komputer Anda. Cek melalui terminal:
```bash
go version
```

### Langkah 2: Buat Folder Proyek Baru
```bash
mkdir proyek-baru
cd proyek-baru
```

### Langkah 3: Inisialisasi Modul Go
```bash
go mod init proyek-baru
```

### Langkah 4: Install Pustaka / Dependensi Eksternal
1. **Gin Gonic (Web Server)**: `go get -u "github.com/gin-gonic/gin"`
2. **PostgreSQL Driver**: `go get -u "github.com/lib/pq"`
3. **SQL Migrate**: `go get -u "github.com/rubenv/sql-migrate"`
4. **Godotenv**: `go get -u "github.com/joho/godotenv"`

---

## 📘 4. Penjelasan Rinci Kode Berdasarkan Urutan Alur (`Entity` ➔ `Migration` ➔ `Database` ➔ `Repository` ➔ `Controller` ➔ `Router` ➔ `Main`)

---

### 4.1. `entities/card_bridge_entity.go` (Wadah Data Kartu)

**Tugas File:** Mendefinisikan bentuk data kartu yang dipakai dalam program Go dan format teks JSON saat dikirim ke pemanggil API.

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```

#### Penjelasan Spesifik Baris demi Baris:
* ◦ `package entities` : Kelompok file ini.
* ◦ `type Card struct` : Membuat cetakan data bernama `Card`.
* ◦ `ID string` : Kotak penyimpanan data bernama `ID` dengan jenis teks (`string`).
* ◦ `json: "id"` : *Struct tag* untuk mengatur nama field menjadi `id` saat dikirim ke format JSON.

---

### 4.2. `databases/sql_migrations/1_initiate.sql` (Skrip Pembuat Tabel)

**Tugas File:** Membuat tabel `card` di PostgreSQL secara otomatis.

```sql
-- +migrate Up
-- +migrate StatementBegin
CREATE TABLE card (id varchar(20));
-- +migrate StatementEnd
```

#### Penjelasan Spesifik:
* ◦ `-- +migrate Up` : Perintah untuk library migrasi agar menjalankan kode ini saat membangun tabel.
* ◦ `CREATE TABLE card` : Perintah SQL untuk membuat tabel baru di database.

---

### 4.3. `databases/database.go` (Pengelola Koneksi & Migrasi Otomatis)

**Tugas File:** Mengurus koneksi ke PostgreSQL dan menjalankan migrasi file `.sql`.

```go
package database
// ... (imports)

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

#### Penjelasan Spesifik:
* ◦ `dbParam *sql.DB` : **Sangat Penting**. Ini adalah *bahan masukan* berupa kunci koneksi database yang dikirim dari `main.go`. Tanpa parameter ini, fungsi migrasi tidak akan tahu ke mana harus mengirim perintah SQL.
* ◦ `n` : Variabel yang menampung *angka integer* jumlah file migrasi yang sukses dijalankan.
* ◦ `errs` : Variabel yang menampung *pesan error* jika terjadi kesalahan (misalnya koneksi database salah).
* ◦ `migrate.Exec(...)` : Fungsi dari library migrasi untuk mengeksekusi perintah SQL.
* ◦ `panic(errs)` : Jika `errs` berisi pesan error, program akan langsung berhenti paksa (`panic`) agar kita sadar ada kesalahan konfigurasi.
* ◦ `DbConnection = dbParam` : Menyimpan kunci koneksi database ke variabel global agar bisa dipakai oleh file lain.

---

### 4.4. `repositories/card_bridge_repo.go` (Akses Data)

**Tugas File:** Tempat koki menjalankan perintah SQL ke database PostgreSQL.

```go
package repositories
// ... (imports)

func GetCards(db *sql.DB) (result []entities.Card, err error) {
	sql := "SELECT * FROM card"
	rows, err := db.Query(sql) // Membaca data
	if err != nil { return }
	defer rows.Close() // Menutup koneksi rows secara aman setelah selesai

	for rows.Next() { // Perulangan membaca baris
		var data entities.Card
		err = rows.Scan(&data.ID) // Memindahkan data dari database ke variabel Go
		if err != nil { return }
		result = append(result, data)
	}
	return
}

func InsertCard(db *sql.DB, card entities.Card) (err error) {
	sql := "INSERT INTO card(id) values($1)"
	_, err = db.Exec(sql, card.ID) // Exec untuk perintah yang tidak mengembalikan data
	return err
}

// ... (UpdateCard & DeleteCard serupa dengan InsertCard)
```

#### Penjelasan Spesifik:
* ◦ `db *sql.DB` : Parameter kunci akses database yang dikirim dari *Controller*.
* ◦ `db.Query(sql)` : Khusus untuk perintah `SELECT` (pembacaan data).
* ◦ `db.Exec(sql, ...)` : Khusus untuk perintah `INSERT`, `UPDATE`, `DELETE` (perubahan data).
* ◦ `rows.Scan(&data.ID)` : Perintah krusial untuk "menarik" isi kolom dari database dan memasukkannya ke dalam alamat memori variabel `data.ID`.

---

### 4.5. `controllers/card_bridge_controllers.go` (Logika Bisnis & HTTP Handler)

**Tugas File:** Pelayan yang menerima permintaan HTTP, memanggil koki (repo), dan mengantarkan makanan (JSON).

```go
// ... (imports)

func GetCards(c *gin.Context) {
	card, err := repositories.GetCards(database.DbConnection) // Panggil koki

	if err != nil { // Jika gagal
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": card}) // Jika sukses
}

// ... (Insert, Update, Delete serupa)
```

#### Penjelasan Spesifik:
* ◦ `c *gin.Context` : Objek konteks Gin. Berisi data request pengguna dan alat untuk mengirim balasan.
* ◦ `c.JSON(http.StatusOK, ...)` : Perintah mengirim balik data dalam format JSON dengan status sukses (`200 OK`).
* ◦ `http.StatusInternalServerError` : Status code `500`, menandakan ada error di dalam server (biasanya error database).
* ◦ `database.DbConnection` : Mengambil kunci koneksi database global untuk dikirim ke Repository.

---

### 4.6. `routers/card_bridge_routers.go` (Pendaftaran Alamat URL)

**Tugas File:** Buku daftar alamat URL dan middleware CORS.

```go
// ... (imports)

func StartServer() *gin.Engine {
	router := gin.Default() // Membuat mesin router Gin

	// Middleware CORS: Satpam yang mengecek apakah domain luar boleh mengakses API
	router.Use(cors.New(cors.Config{...})) 
	
	router.GET("/cards", controllers.GetCards)
	router.POST("/card/input/:id", controllers.InsertCard)
	router.PUT("/card/update/:id", controllers.UpdateCard)
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	return router
}
```

#### Penjelasan Spesifik:
* ◦ `gin.Default()` : Mengaktifkan mesin web server Gin dengan fitur *Logger* (pencatat aktivitas di terminal) dan *Recovery* (penahan *crash*).
* ◦ `cors.Config` : Pengaturan yang mengizinkan domain mana saja yang boleh mengakses API ini (`AllowOrigins: "*"`) dan metode HTTP apa saja yang diizinkan (`AllowMethods`).

---

### 4.7. `main.go` (Pintu Masuk Aplikasi)

**Tugas File:** Program utama untuk menyalakan koneksi database, migrasi, dan server web.

```go
// ... (imports)

func main() {
	var PORT = ":50000" // Port server

	// Rangkaian string koneksi
	psqlInfo := fmt.Sprintf(...)

	DB, err = sql.Open("postgres", psqlInfo) // Buka pintu database
	if err != nil { log.Fatal(...) }

	database.DBMigration(DB) // Panggil fungsi migrasi otomatis
	defer DB.Close() // Pastikan database ditutup saat aplikasi mati

	routers.StartServer().Run(PORT) // Jalankan web server
}
```

#### Penjelasan Spesifik:
* ◦ `sql.Open("postgres", psqlInfo)` : Perintah untuk membuka gerbang koneksi ke database.
* ◦ `database.DBMigration(DB)` : Memanggil fungsi migrasi di folder `databases` dengan mengirimkan kunci koneksi `DB`.
* ◦ `defer DB.Close()` : Perintah jaminan agar saat aplikasi ditutup (`Ctrl+C`), koneksi database diputus dengan bersih.
* ◦ `routers.StartServer().Run(PORT)` : Perintah pamungkas untuk menyalakan web server HTTP di port yang ditentukan. Fungsi `Run` bersifat *blocking*, artinya program tidak akan lanjut ke baris berikutnya sampai server dihentikan.
