$API_URL = "http://localhost:8080/api/tasks"

# Loop from 1 to 20
for ($i = 1; $i -le 20; $i++) {
    $Title = "Task No. $i - Pagination Test"

    $body = @{
        title = $Title
    }

    Write-Host "Sending task No. $i..."


    Invoke-RestMethod -Method Post `
        -Uri $API_URL `
        -Body ($body | ConvertTo-Json) `
        -ContentType "application/json"
}

Write-Host "successfully!"