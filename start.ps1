function BuildProto {
    Write-Host "Compiling from gRPC proto files"
    Remove-Item -Path "src/proto/*" -Recurse -Force
    protoc --go_out=. --go-grpc_out=. proto/comm.proto
}

function StartServer {
    param (
        [string]$Addr,
        [int]$Port = 60000,
        [int]$Timeout = 5
    )
    if (-not $Addr) {
        $Addr = Get-WiFiIPAddress
        if (-not $Addr) {
            $Addr = "0.0.0.0"
        }
    }

    Write-Host "Starting node on ${Addr}:${Port}"
    Start-Process -FilePath "cmd" -ArgumentList "/c start cmd /k `"go run src/node/main.go -addr ${Addr}:${Port} -timeout $Timeout`"" -NoNewWindow
}

function StartServers {
    param (
        [int]$Size = 1,
        [string]$Addr,
        [int]$Port = 60000,
        [int]$Timeout = 5,
        [string]$Hostfile = "./config/Hostfile"
    )
    if (-not $Addr) {
        $Addr = Get-WiFiIPAddress
        if (-not $Addr) {
            $Addr = "0.0.0.0"
        }
    }

    Write-Host "Generating hostfile"
    Remove-Item -Path $Hostfile -Recurse -Force
    
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        $address = "${Addr}:${nPort}"
        $address | Out-File -FilePath $Hostfile -Append -Encoding ascii
    }

    Write-Host "Starting $Size servers"
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        StartServer -Addr $nAddr -Port $nPort -Timeout $Timeout -Hostfile $Hostfile
    }

    Controller
}

function Controller {
    go run src/controller/main.go
}


function Get-WiFiIPAddress {
    $wifiIPAddress = Get-NetIPAddress | Where-Object -FilterScript { $_.InterfaceAlias -Eq "Wi-Fi" -and $_.AddressFamily -Eq "IPv4" }
    $ipAddress = $wifiIPAddress.IPAddress
    return $ipAddress
}