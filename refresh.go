package main

import(
   "io"
   "os"
   "fmt"
   "strings"
   "net/http"
   "path/filepath"
)

// 移除危險字符和路徑遍歷
func sanitizeFileName(filename string)(string) {
   filename = strings.ReplaceAll(filename, "..", "")
   filename = strings.ReplaceAll(filename, "/", "")
   filename = strings.ReplaceAll(filename, "\\", "")
   filename = strings.TrimSpace(filename)
   if len(filename) == 0 {
      return ""
   }
   return filename
}

func refreshScreen(w http.ResponseWriter, r *http.Request) {
   if r.Method != "POST" {
      http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
      return
   }
   // 建立目標檔案
   uploadDir := os.Getenv("uploadDir")
   if uploadDir != "" {
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprintf(w, `<div class="message error">❌ 系統參數upload dir設定錯誤</div>`)
      return
   }
   // 確認目錄是否存在，不存在則建立
   if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
      http.Error(w, "Unable to create directory", http.StatusInternalServerError)
      return
   }
   // 處理上傳檔案
   r.ParseMultipartForm(32 << 20) // 32MB
   file, handler, err := r.FormFile("file")
   if err != nil {
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprintf(w, `<div class="message error">❌ 錯誤：%s</div>`, err.Error())
      return
   }
   defer file.Close()
   filename := sanitizeFileName(handler.Filename)
   if filename == "" {
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprintf(w, `<div class="message error">❌ 檔案名稱無效</div>`)
      return
   }
   // 複製檔案內容到目標檔案 layerssid.bmp
   destFileName := os.Getenv("screenFileName")
   if destFileName == "" {
      w.WriteHeader(http.StatusBadRequest)
      fmt.Fprintf(w, `<div class="message error">❌ 系統參數screen file name 設定錯誤</div>`)
      return
   }
   destPath := filepath.Join(uploadDir, destFileName)
   dst, err := os.Create(destPath) // os.Create() 如果檔案已存在會直接覆蓋掉
   if err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, `<div class="message error">❌ 建立檔案失敗：%s</div>`, err.Error())
      return
   }
   defer dst.Close()
   if _, err := io.Copy(dst, file); err != nil {
      w.WriteHeader(http.StatusInternalServerError)
      fmt.Fprintf(w, `<div class="message error">❌ 儲存檔案失敗：%s</div>`, err.Error())
      return
   }
   fmt.Fprintf(w, `<div class="message success">✅ 檔案 "%s" 上傳成功！</div>`, filename)
   go afterUpload(destPath, filepath.Join(uploadDir, "layerssid.bmp"))
}
