function BuildProto {
    Write-Host "Compiling from gRPC proto files"
    Remove-Item -Path "src/proto/*" -Recurse -Force
    protoc --go_out=. --go-grpc_out=. proto/comm.proto
}

function StartServer {
    param (
        [string]$Addr = "0.0.0.0:60000",
        [int]$Timeout = 5
    )
    Write-Host "Starting node on $Addr"
    Start-Process -FilePath "cmd" -ArgumentList "/c start cmd /k `"go run src/main.go -addr $Addr -timeout $Timeout`"" -NoNewWindow
}

function StartServers {
    param (
        [int]$Size = 1,
        [string]$Addr = "0.0.0.0",
        [int]$Port = 60000,
        [int]$Timeout = 5,
        [string]$Hostfile = "./config/Hostfile"
    )

    Write-Host "Generating hostfile"
    Remove-Item -Path $Hostfile -Recurse -Force
    
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        $address = "0.0.0.0:$nPort"
        $address | Out-File -FilePath $Hostfile -Append -Encoding ascii
    }

    Write-Host "Starting $Size servers"
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        $nAddr = "${Addr}:${nPort}"
        StartServer -Addr $nAddr -Timeout $Timeout -Hostfile $Hostfile
    }   
}

function Get-WiFiIPAddress {
    $wifiInterface = Get-NetAdapter | Where-Object {$_.InterfaceDescription -like "*Wireless*"}
    $wifiIPAddress = Get-NetIPAddress -InterfaceIndex $wifiInterface.ifIndex | Where-Object {$_.AddressFamily -eq "IPv4"} | Select-Object -ExpandProperty IPAddress
    return $wifiIPAddress
}