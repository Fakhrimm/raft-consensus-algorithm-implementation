function BuildProto {
    Write-Host "Compiling from gRPC proto files"
    Remove-Item -Path "src/proto/*" -Recurse -Force
    protoc --go_out=. --go-grpc_out=. proto/comm.proto
}

function StartServer {
    param (
        [int]$Port = 60000
    )
    Write-Host "Starting node on port $Port"
    Start-Process -FilePath "cmd" -ArgumentList "/c start cmd /k `"go run src/main.go -port $Port`"" -NoNewWindow
}

function StartServers {
    param (
        [int]$Size = 1,
        [int]$Port = 60000,
        [string]$Hostfile = "./config/Hostfile"
    )

    Write-Host "Generating hostfile"

    "" | Out-File -FilePath $Hostfile -Force
    
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        $address = "0.0.0.0:$nPort"
        $address | Out-File -FilePath $Hostfile -Append
    }

    Write-Host "Starting $Size servers"
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        StartServer -Port $nPort -Hostfile $Hostfile
    }   
}