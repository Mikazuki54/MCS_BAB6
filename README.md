# Panduan & Dokumentasi Kode Lengkap REST API Card Bridge (`bab6`)
*Panduan Ramah Pemula & Orang Awam: Penjelasan Kode 100% Mendalam, Bedah Simbol/Sintaks, serta Bedah Tuntas Setiap Baris Kode di Dalam Setiap Function.*

---

## 💡 1. Konsep Dasar & Analogia Sederhana (Untuk Orang Awam)

Analogi: **Restoran**.
- **ROUTER**: Kasir yang mengarahkan pesanan.
- **CONTROLLER**: Pelayan yang memastikan format pesanan benar dan menyajikan makanan (JSON).
- **REPOSITORY**: Koki yang memasak/mengambil bahan langsung ke gudang (Database).
- **DATABASE**: PostgreSQL (Gudang data).

---

## 📖 2. Glosarium Simbol & Istilah Penting

| Simbol / Istilah | Penjelasan Sederhana |
| :--- | :--- |
| **`//go:embed`** | Perintah untuk membungkus file eksternal (seperti `.sql`) ke dalam binary aplikasi. |
| **`_` (Blank Identifier)** | Variabel "tempat sampah" untuk mengabaikan hasil fungsi yang tidak dipakai agar tidak error. |
| **`*` (Pointer)** | Kunci akses langsung ke alamat memori asli data. |
| **`&` (Address-of)** | Mengambil lokasi memori variabel agar bisa diisi nilainya oleh fungsi lain (misal dari database). |
| **`:=`** | Cara cepat membuat & mengisi variabel baru. |
| **`defer`** | Perintah "tunda dulu" agar kode dijalankan saat fungsi hampir selesai. |
| **CORS** | Keamanan browser agar website luar tidak sembarangan mengambil data API kita. |

---

## 🛠️ 3. Langkah-Langkah Membuat Proyek Go Baru

1. `mkdir proyek-baru` & `cd proyek-baru`
2. `go mod init proyek-baru`
3. `go get github.com/gin-gonic/gin`
4. `go get github.com/lib/pq`
5. `go get github.com/rubenv/sql-migrate`

---

## 📘 4. Bedah Tuntas Setiap Function per File

---

### 4.1. `entities/card_bridge_entity.go`
**Tugas:** Menentukan wadah cetakan data (`struct`).

```go
package entities

type Card struct {
	ID string `json: "id"`
}
```
* **Penjelasan:**
  * `type Card struct`: Membuat cetakan data bernama `Card`.
  * `ID string`: Wadah teks bernama `ID`.
  * `` `json: "id"` ``: Agar saat dikirim ke format JSON, nama fieldnya menjadi `id`.

---

### 4.2. `databases/sql_migrations/1_initiate.sql`
**Tugas:** Skrip pembuat tabel di database.

```sql
-- +migrate Up
CREATE TABLE card (id varchar(20));
```
* **Penjelasan:** `-- +migrate Up` adalah penanda untuk alat migrasi agar menjalankan perintah SQL di bawahnya.

---

### 4.3. `databases/database.go`
**Tugas:** Mengurus koneksi database dan menjalankan migrasi.

```go
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
* **Bedah Fungsi `DBMigration(dbParam *sql.DB)`**:
  * `dbParam *sql.DB`: Kunci koneksi database yang diterima dari `main.go`.
  * `migrations := ...`: Menyiapkan konfigurasi file SQL dari hasil *embed*.
  * `migrate.Exec(...)`: Perintah utama menjalankan migrasi.
    * `n`: Menyimpan jumlah file SQL yang sukses.
    * `errs`: Menyimpan pesan error jika gagal.
  * `if errs != nil { panic(errs) }`: Jika ada error, hentikan program paksa.
  * `DbConnection = dbParam`: Simpan koneksi ke variabel global agar bisa dipakai file lain.

---

### 4.4. `repositories/card_bridge_repo.go`
**Tugas:** Tempat menjalankan perintah SQL ke database.

#### A. Fungsi `GetCards(db *sql.DB) (result []entities.Card, err error)`
```go
func GetCards(db *sql.DB) (result []entities.Card, err error) {
	sql := "SELECT * FROM card"            // 1. Siapkan string query SQL
	rows, err := db.Query(sql)             // 2. Eksekusi query SELECT
	if err != nil { return }               // 3. Jika error, keluar fungsi
	defer rows.Close()                     // 4. Pastikan rows ditutup di akhir

	for rows.Next() {                      // 5. Perulangan membaca hasil query
		var data entities.Card             // 6. Siapkan wadah 1 data kartu
		err = rows.Scan(&data.ID)          // 7. Salin kolom database ke memori &data.ID
		if err != nil { return }           // 8. Cek error saat scan
		result = append(result, data)      // 9. Masukkan ke daftar hasil
	}
	return                                 // 10. Kembalikan data dan error
}
```

#### B. Fungsi `InsertCard(db *sql.DB, card entities.Card) (err error)`
```go
func InsertCard(db *sql.DB, card entities.Card) (err error) {
	sql := "INSERT INTO card(id) values($1)" // 1. Siapkan query INSERT
	_, err = db.Exec(sql, card.ID)           // 2. Eksekusi perintah (ignore metadata dengan _)
	return err                               // 3. Kembalikan error
}
```

#### C. Fungsi `UpdateCard` & `DeleteCard`
* Isinya sama dengan `InsertCard`, hanya berbeda perintah SQL (`UPDATE` atau `DELETE`) dan jumlah parameter.

---

### 4.5. `controllers/card_bridge_controllers.go`
**Tugas:** Menerima permintaan HTTP, memanggil Repository, dan mengantar balasan JSON.

#### A. Fungsi `GetCards(c *gin.Context)`
```go
func GetCards(c *gin.Context) {
	var result gin.H
	card, err := repositories.GetCards(database.DbConnection) // 1. Panggil repo
	if err != nil { result = gin.H{"result": err.Error()} } // 2. Jika gagal
	else { result = gin.H{"result": card} }                   // 3. Jika sukses
	c.JSON(http.StatusOK, result)                             // 4. Kirim JSON
}
```

#### B. Fungsi `InsertCard`, `UpdateCard`, `DeleteCard`
* Intinya sama: `c.Param` atau `c.ShouldBindJSON` mengambil data dari user -> Panggil `repositories` -> Jika `err != nil` kirim error, jika sukses kirim pesan berhasil dengan `c.JSON`.

---

### 4.6. `routers/card_bridge_routers.go`
**Tugas:** Mendaftarkan alamat URL ke fungsi di Controller.

```go
func StartServer() *gin.Engine {
	router := gin.Default() // 1. Buat engine router
	
    // 2. Middleware CORS (Satpam Browser)
	router.Use(cors.New(cors.Config{...})) 
	
	router.GET("/cards", controllers.GetCards) // 3. Daftar rute URL
	// ... (rute lainnya)

	return router // 4. Kembalikan engine
}
```

---

### 4.7. `main.go` (Entry Point)
**Tugas:** Menghubungkan database, migrasi, dan menyalakan server web.

```go
func main() {
	var PORT = ":50000"
    // 1. Konfigurasi String PostgreSQL
	psqlInfo := fmt.Sprintf(...)

    // 2. Buka koneksi database
	DB, err = sql.Open("postgres", psqlInfo)
	if err != nil { log.Fatal(...) }

    // 3. Jalankan migrasi
	database.DBMigration(DB) 
	defer DB.Close() // 4. Tutup koneksi saat program mati

    // 5. Jalankan server web
	routers.StartServer().Run(PORT) 
}
```
