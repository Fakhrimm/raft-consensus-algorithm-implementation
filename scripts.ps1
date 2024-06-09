function Clean {
    Get-ChildItem -Path "config" -Recurse -Filter *.storage | Remove-Item -Force
}

function Tidy {
    Write-Host "Running 'go mod tidy' at node and controller"

    Push-Location -Path "node"
    go mod tidy
    Pop-Location

    Push-Location -Path "controller"
    go mod tidy
    Pop-Location
}

function Proto {
    Write-Host "Compiling from gRPC proto files"
    Remove-Item -Path "grpc/*" -Recurse -Force
    protoc --go_out=. --go-grpc_out=. proto/comm.proto

    $destDirs = @("node/grpc", "controller/grpc")

    foreach ($destDir in $destDirs) {
        if (-Not (Test-Path $destDir)) {
            New-Item -ItemType Directory -Force -Path $destDir
        }
        Copy-Item -Path "grpc/*" -Destination $destDir -Recurse -Force
    }
    Remove-Item -Path "grpc" -Recurse -Force

    Push-Location -Path "client"
    Remove-Item -Path "./src/grpc" -Recurse -Force
    New-Item -ItemType Directory -Force -Path "./src/grpc" > $null
    npx protoc --ts_out "./src/grpc" --proto_path ../proto ../proto/comm.proto 
    Pop-Location

    Write-Host "Proto files compiled and copied successfully."
}

function Client {
    Push-Location -Path "client"
    npm run dev
    Pop-Location
}

function Server {
    param (
        [string]$Addr,
        [int]$Port = 60000,
        [int]$Timeout = 2000,
        [boolean]$IsJointConsensus = $false,
        [string]$Hostfile = "/config/Hostfile",
        [string]$HostfileNew = "default"
    )
    if (-not $Addr) {
        $Addr = WifiIP
        if (-not $Addr) {
            $Addr = "0.0.0.0"
        }
    }
    Write-Host "Starting node on ${Addr}:${Port}"
    
    Push-Location -Path "node"
    $command = "go run main.go -addr ${Addr}:${Port} -hostfile $Hostfile -hostfilenew $HostfileNew -timeout $Timeout"

    if ($IsJointConsensus) {
        $command += " -isjointconsensus"
    }

    # Write-Host $command
    if (1 -eq 2) {
        Start-Process -FilePath "cmd" -ArgumentList "/c start cmd /k", $command -NoNewWindow > $null
    }
    else {
        StartWithWindowsTerminal -Command $command "${Addr}:${Port}" > $null
    }
    Pop-Location
}

function Servers {
    param (
        [int]$Size = 1,
        [string]$Addr,
        [int]$Port = 60000,
        [int]$Timeout = 2000,
        [string]$Hostfile = "/config/Hostfile"
    )
    if (-not $Addr) {
        $Addr = WifiIP
        if (-not $Addr) {
            $Addr = "0.0.0.0"
        }
    }

    Write-Host "Generating hostfile"
    Remove-Item -Path ".$Hostfile" -Recurse -Force
    
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        $address = "${Addr}:${nPort}"
        $address | Out-File -FilePath ".$Hostfile" -Append -Encoding ascii
    }

    Write-Host "Starting $Size servers"
    for ($i = 0; $i -lt $Size; $i++) {
        $nPort = $Port + $i
        Server -Addr $nAddr -Port $nPort -Timeout $Timeout -Hostfile $Hostfile
    }

    Controller
}

function Controller {
    Push-Location -Path "controller"
    go run main.go
    Pop-Location
}


function WifiIP {
    $wifiIPAddress = Get-NetIPAddress | Where-Object -FilterScript { $_.InterfaceAlias -Eq "Wi-Fi" -and $_.AddressFamily -Eq "IPv4" }
    $ipAddress = $wifiIPAddress.IPAddress
    return $ipAddress
}

function StartWithWindowsTerminal {
    param (
        [string]$Command,
        [string]$Title = "Server"
    )
    $initDir = Get-Location
    $startInfo = New-Object System.Diagnostics.ProcessStartInfo
    $startInfo.FileName = "wt.exe"
    $startInfo.Arguments = "-w 0 nt --suppressApplicationTitle --title " + $Title + " -d " + $initDir + " cmd /k " + $Command
    $startInfo.UseShellExecute = $false

    #     Set environment variables (Uncomment if needed)
    #    $startInfo.Environment["Path"] = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
    #            [System.Environment]::GetEnvironmentVariable("Path", "User")

    # Start the process
    [System.Diagnostics.Process]::Start($startInfo)
}