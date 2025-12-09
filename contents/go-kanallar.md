*"Don't communicate by sharing memory, share memory by communicating."*

Go-nu populyarlığa qaldıran ən önəmli faktorlardan biri onun **concurrency modeli** idi. Ənənəvi dillərdə istifadə olunan **mutex**, **lock**, **semaphore** kimi primitivlər çox zaman kodun yavaş, qarışıq və deadlocklarla dolu olmasına gətirib çıxarırdı. Go isə bu problemi **channel** əsaslı yanaşma ilə tam fərqli şəkildə həll edir.

## Channel nədir?

Channel, Go-da **goroutinlər arasında məlumat ötürmək üçün portaldır**. Bir goroutine kanala dəyər göndərir, digəri isə həmin dəyəri qəbul edir. Bu yanaşma sayəsində paylaşılmış yaddaşla uğraşmadan sinxron kommunikasiya mümkündür. Kanallarda məlumatın göndərilməsi və alınması eyni anda baş verməlidir. Eyni bir portal kimi məlumat bir nöqtədən girib digərindən çıxmalıdır.

Sadə nümunə:

```go
ch := make(chan int)

go func() {
    ch <- 42 // göndər
}()

value := <-ch // qəbul et
fmt.Println(value)
```

## Buffered vs Unbuffered Channels

Go-da iki növ channel mövcuddur və onların davranışı bir-birindən ciddi fərqlənir:

**1. Unbuffered Channel**

Unbuffered channel yaratmaq üçün sadəcə `make(chan Type)` yazırıq:

```go
ch := make(chan int) // unbuffered channel
```

Unbuffered channellarda göndərən tərəf **bloklanır** və qəbul edən tərəf hazır olana qədər gözləyir. 

```go
ch := make(chan int)

go func() {
    fmt.Println("Göndərməyə hazıram...")
    ch <- 100 // bu sətirdə bloklanacaq
    fmt.Println("Göndərdim!") // yalnız qarşı tərəf qəbul edəndən sonra
}()

time.Sleep(2 * time.Second) 
value := <-ch // indi göndərən goroutine davam edə bilər
fmt.Println("Aldım:", value)
```

**2. Buffered Channel**

Buffered channel yaratmaq üçün isə ikinci parametr olaraq bufer ölçüsü göstəririk:

```go
ch := make(chan int, 3) // 3 ölçülü buferli channel
```

Buffered channellarda göndərən yalnız **bufer dolduqda** bloklanır. Buferdə yer varsa, göndərmə dərhal baş verir:

```go
ch := make(chan int, 3)

ch <- 1 // bloklanmır
ch <- 2 // bloklanmır  
ch <- 3 // hələ bloklanmır, bufer doludur amma yerləşir

ch <- 4 // BURADA bloklanacaq, çünki bufer artıq tamamilə doludur

value := <-ch // birini çıxartdıq, indi yenə yer var
ch <- 5 // bu rahat yerləşəcək
```

## Channel-ı bağlamaq və Range istifadəsi

Channel-ı bağlamaq üçün `close()` funksiyasından istifadə edirik. Bu çox vaxt göndərən tərəf tərəfindən edilir:

```go
ch := make(chan int, 5)

// göndərən goroutine
go func() {
    for i := 1; i <= 5; i++ {
        ch <- i
        fmt.Println("Göndərildi:", i)
    }
    close(ch) // işim bitdi, kanalı bağlayıram
}()

// qəbul edən tərəf
for value := range ch {
    fmt.Println("Alındı:", value)
}
// kanal bağlananda range döngüsü avtomatik bitir
```

Bağlanmış kanala yenidən göndərmə etmək **panic** yaradır. Amma bağlanmış kanaldan oxumaq təhlükəsizdir - sıfır dəyər və `false` qaytarır:

```go
ch := make(chan int)
close(ch)

value, ok := <-ch
fmt.Println(value, ok) // 0 false
```

## Select - çox kanallı gözləmə

`select` statement-i bir neçə channel əməliyyatından birinin hazır olmasını gözləməyə imkan verir. İlk hansı channel hazır olarsa, həmin case işə düşür:

```go
ch1 := make(chan string)
ch2 := make(chan string)

go func() {
    time.Sleep(1 * time.Second)
    ch1 <- "birinci kanaldan"
}()

go func() {
    time.Sleep(2 * time.Second)
    ch2 <- "ikinci kanaldan"
}()

for i := 0; i < 2; i++ {
    select {
    case msg1 := <-ch1:
        fmt.Println("ch1:", msg1)
    case msg2 := <-ch2:
        fmt.Println("ch2:", msg2)
    }
}
```

Select-ə **timeout** və **default** case əlavə etmək də mümkündür:

```go
select {
case msg := <-ch:
    fmt.Println("Mesaj alındı:", msg)
case <-time.After(3 * time.Second):
    fmt.Println("Timeout! 3 saniyə gözlədim, heç nə gəlmədi")
default:
    fmt.Println("Hazırda hazır olan heç nə yoxdur")
}
```

Default case varsa, select heç vaxt bloklanmır - hazır channel yoxdursa dərhal default case işləyir.