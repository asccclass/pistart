package main

import (
   "os"
   "fmt"
   "net"
   "time"
   "encoding/json"
)

// Info 結構用於存放需輸出的資訊
type Info struct {
   MACAddress string `json:"mac_address"`
   IPAddress  string `json:"ip_address"`
   SystemTime string `json:"system_time"`
}

func getMACAndIP() (macAddr, ipAddr string) {
   interfaces, err := net.Interfaces()
   if err != nil {
      return "", ""
   }

   for _, iface := range interfaces {
      // 跳過 loopback 和沒有啟用的介面
      if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
         continue
      }

      // 取得 MAC Address
      if mac := iface.HardwareAddr.String(); mac != "" {
         macAddr = mac
      }

      addrs, err := iface.Addrs()
      if err != nil {
         continue
      }

      for _, addr := range addrs {
         var ip net.IP
         switch v := addr.(type) {
         case *net.IPNet:
            ip = v.IP
         case *net.IPAddr:
            ip = v.IP
         }
         // 取得 IPv4 Address
         if ip != nil && ip.To4() != nil {
            ipAddr = ip.String()
            // 只要找到一組就跳出
            return macAddr, ipAddr
         }
      }
   }
   return macAddr, ipAddr
}

func main() {
   mac, ip := getMACAndIP()
   // 取得當前系統時間並格式化
   now := time.Now().Format("2006-01-02 15:04:05")

   info := Info{
      MACAddress: mac,
      IPAddress:  ip,
      SystemTime: now,
   }

   // 輸出 JSON 格式
   jb, err := json.MarshalIndent(info, "", "  ")
   if err != nil {
      fmt.Fprintf(os.Stderr, "JSON encode error: %v\n", err)
      os.Exit(1)
   }
   fmt.Println(string(jb))
}
