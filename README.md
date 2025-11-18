# Практическая работа 4   
## Выполнил Сотников М. ЭФМО-01-25

Этот проект представляет собой простой RESTful API для управления задачами (ToDo List), разработанный на языке Go с использованием популярного маршрутизатора go-chi/chi. Стандартный роутер Go, `http.ServeMux`, очень минималистичен: он маршрутизирует запросы только по сопоставлению префиксов (началу URL). Он не имеет встроенной поддержки для разных HTTP-методов (GET, POST) и не позволяет легко извлекать переменные из пути (например, ID из /tasks/123).

Специализированные роутеры, такие как `go-chi/chi`, предоставляют полную поддержку шаблонов, позволяя определять маршруты вроде /tasks/{id} и автоматически извлекать id.

Реализованы CRUD операции (POST, GET, PUT, DELETE).  


<img width="855" height="375" alt="health200" src="https://github.com/user-attachments/assets/b30265ec-257a-4e81-8c2f-18bb8bc187ac" />    


<img width="866" height="779" alt="PostPostman" src="https://github.com/user-attachments/assets/128ba622-e1f2-4250-92ec-12693d558533" />    


<img width="838" height="731" alt="GetTasks" src="https://github.com/user-attachments/assets/afa4ffe4-e940-443d-954b-0f7078ad5fcf" />


<img width="827" height="473" alt="GetTaskWithId" src="https://github.com/user-attachments/assets/b6e97c64-665e-4599-90a8-6f11d966cc96" />    


<img width="842" height="494" alt="PutTaskWithId" src="https://github.com/user-attachments/assets/08db4904-b62b-40d3-b112-f8cc9a8bd9f6" />  


<img width="846" height="234" alt="DeleteTask" src="https://github.com/user-attachments/assets/3ebaaaad-4b8b-4747-b4db-db6b8fc3298e" />    


<img width="834" height="550" alt="DeleteTaskChech" src="https://github.com/user-attachments/assets/81236f78-2d05-48a1-aff5-ebebd94da107" />    


<img width="437" height="204" alt="Logger" src="https://github.com/user-attachments/assets/3814552a-8891-40bb-992c-d07a9ec7b081" />     



**Middleware** - функция, которая оборачивает основной обработчик запроса и позволяет выполнить определенную логику до или после того, как запрос достигнет конечного обработчика.

Стандартные middleware **RequestID**(идентификация запроса) и **Recoverer**(безопасное восстановление после сбоев) прописаны в файле main.go и импортируются из пакета github.com/go-chi/chi/v5/middleware под псевдонимом chimw.

Кастомный middleware, находятся внутри пакета pkg/middleware.

**Logger**  
Этот middleware отвечает за логирование всех входящих запросов. Он регистрирует метод, путь, время выполнения и статус ответа.

**SimpleCORS**
Этот middleware отвечает за добавление необходимых HTTP-заголовков (Access-Control-Allow-Origin, Access-Control-Allow-Methods и т.д.), которые позволяют веб-браузерам выполнять междоменные запросы к вашему API.




## Валидация длины заголовка

```golang
const ( //константы валидации длины title
	minTitleLen = 3
	maxTitleLen = 100
)

```    

```golang
titleLen := len(req.Title)
	if titleLen < minTitleLen || titleLen > maxTitleLen {
		httpError(w, http.StatusBadRequest, "title length must be between 3 and 100 characters")
		return
	}
```
<img width="1386" height="586" alt="LengthVal" src="https://github.com/user-attachments/assets/671f9dfc-60dd-4f90-ab0a-35befb335d14" />

### Пагинация
Скрипт для powershell для заполнения задачами:
```bash
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
```
<img width="553" height="664" alt="Success20Posts" src="https://github.com/user-attachments/assets/9b3714c8-bc3a-4332-a03c-1a314c6df224" />    

<img width="574" height="947" alt="DefaultPagination10" src="https://github.com/user-attachments/assets/1d589f3b-a612-46c6-b390-47ea293d3f8d" />    

<img width="581" height="622" alt="Page2Limit5" src="https://github.com/user-attachments/assets/d95bffb5-2a5f-4924-a8fd-e5aad758259f" />   









