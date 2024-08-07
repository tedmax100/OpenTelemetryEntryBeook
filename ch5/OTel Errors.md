---
tags:
  - open-telemetry
---
# OTel Errors

<https://opentelemetry.io/blog/2024/otel-errors/>



當不同的語言對於錯誤（Error）或異常（Exception）的定義以及處理方式存在差異時，例如 Go 就沒有exception，而是用panic。Go 鼓勵開發者將普通的error視為正常的情況，然後返回普通的錯誤值在處理這些情況，而不是全部的問題都視為無法處理的異常。而其他語言如 Java 則支持拋出和捕獲異常。

因此需要在這些語言中標準化錯誤或異常的定義。OpenTelemetry是解決這個問題的工具之一，它提供了以下功能：

1. **錯誤在後端的可視化：** 後端中錯誤的呈現方式可能與您預期的不同，或者與您期望的不同。OpenTelemetry可以幫助標準化錯誤的呈現方式，使得無論使用哪種語言開發的微服務，錯誤報告都能保持一致。

2. **Span類型對錯誤報告的影響：** Span 類型是指追蹤系統中的一個事件，它可以是一個HTTP請求、一個函數調用或者一個資料庫查詢等。不同類型的 Span 對錯誤報告可能有不同的影響，OpenTelemetry可以幫助您管理這些差異。

3. **Span 報告的錯誤與日誌的錯誤：** OpenTelemetry可以幫助區分Span報告的錯誤和日誌中的錯誤，從而更好地理解和分析系統中的錯誤。

總之，OpenTelemetry是一個可以標準化遙測和錯誤報告的工具，無論使用哪種語言開發的微服務都可以使用它來保持一致的錯誤報告機制。



## **Errors versus exceptions**



在軟體開發中，Error 和 Exception 兩個常見但有所不同的概念。

1. **Error**：Error是程序執行過程中的意外問題，它們可能來自於語法錯誤、邏輯錯誤或者其他未預料到的問題。例如，在程式碼中漏掉了一個分號或者使用了未定義的變數，這都可能導致錯誤。錯誤通常發生在程式碼執行之前（編譯時錯誤）或執行中（運行時錯誤）。

2. **Exception** ：Exception是程序執行過程中的一種特殊情況，它會**中斷**程序的正常流程。當程序遇到無法處理的錯誤情況時，通常會拋出一個異常，並且程序的執行流程會被轉移到異常處理代碼中。異常可以是由程式設計者顯式拋出的（例如，除以零），也可以是由系統自動拋出的（例如，空指針異常）。

一些程式語言將錯誤和異常視為相同的概念，並使用相同的機制來處理它們。其他語言則將它們區分開來，並提供不同的機制來處理錯誤和異常。了解這兩者之間的區別對於有效的錯誤處理至關重要，因為它有助於選擇適當的處理策略並確保程序的穩定性和可靠性。



## **Handling errors in OTel**



OTel 是一個跨語言的標準化工具，用於收集應用程序的遙測數據，包括錯誤資訊。由於不同語言之間存在著錯誤和異常處理的概念差異，因此OTel規範（Spec）提供了一個統一的標準，用於各語言的實現，從而實現跨語言的一致性。



在OTel中，語言API和SDK的實現是基於該規範的，並且有一些通用的規則限制了在實現中可以添加的功能，這保證了各語言的一致性。儘管如此，在某些情況下，一些語言可能會為了特定的原因而偏離規範，並提供了一些彈性以便每種語言盡可能地以該語言的習慣方式來實現功能。在實踐中，也有一些例外情況；例如，某種語言可能會將一項新功能作為將其添加到規範的一部分的原型，但在相應的語言被添加之前，該功能可能會被發布（通常為 alpha 或 experimental）。該版本的狀態於本書的 Ch 4.4 中有介紹。



另一個例外情況是，當一種語言決定偏離規範時。雖然一般不建議這樣做，但有時有強烈的語言特定原因來做一些不同的事情。這樣，規範允許每種語言在實現特性時具有一定的靈活性，使得其能夠盡可能符合該語言的習慣用法。例如，大多數語言都實現了`RecordException`，而Go則實現了`RecordError`，其作用相同。



我們能看到  [OTel Spec 的合規性矩陣](https://github.com/open-telemetry/opentelemetry-specification/blob/main/spec-compliance-matrix.md)中，Go 在Record Exception 就處於尚未支持。但其實用的是 RecordError。

![img][圖1]



### Errors in logs



在 OTel中，日誌是由服務或其他組件發出的結構化、時間戳記的資訊。這個最近添加到OTel中的功能提供了另一種報告錯誤的方式。日誌通常具有不同的嚴重性級別，用於表示消息的類型，例如DEBUG、INFO、WARNING、ERROR和CRITICAL。

OTel允許將日誌與 Span 相關聯，其中一條日誌消息可以通過 Trace Context 關聯到 Span 內的一個 Span，從而將日誌和 Span 關聯起來。因此，查找具有ERROR或CRITICAL日誌級別的日誌資訊可以進一步獲取導致該錯誤的資訊，通過提取相關的 Span 資訊。

要在日誌中記錄錯誤，需要提供 `exception.type` 或 `exception.message`，這是必需的資訊。此外，建議還提供 `exception.stacktrace`，以提供錯誤的堆疊追蹤資訊。有關更多資訊，可以參考日誌中異常的語義慣例（[Semantic Conventions for Exceptions in Logs](https://app.heptabase.com/874af36f-34e9-4f16-9008-ae94afbbfb6f/card/f6bae872-21b5-4702-b3d0-691038bd0df1)）。總的來說，日誌功能在OTel中提供了一個強大的錯誤報告機制，使開發人員能夠更有效地處理和監控運行的應用程序中的錯誤情況。



### Error in Spans



在 OTel 中，Span 是分散式追蹤的基本元件，代表分散式系統中的個別工作單位。Span 通過上下文與其他 Span 及追蹤相關聯。簡單來說，上下文是將一組數據轉換為統一追蹤的粘合劑。上下文傳播允許我們在多個系統之間傳遞資訊，從而將它們聯繫在一起。通過 metadata 和 Span 事件，追蹤可以告訴我們有關運行中的應用程序的各種訊息。



![img][圖2]


在 OTel 中，您可以使用 metadata （屬性）來增強 Span。通過將相關資訊（如用戶 ID、請求參數或環境變數）附加到 Span 上，您可以更深入地了解錯誤發生的情況，並快速確定其根本原因。這種富有metadata的錯誤處理方法可以顯著減少診斷和解決問題所需的時間和精力，從而提高應用程序的可靠性和可維護性。



此外，Span 還具有 Span 類型欄位，它提供了一些額外的metadata，可以幫助開發人員進行錯誤排查。OTel 定義了幾種不同的 Span 類型，每種類型對錯誤報告都有獨特的影響。Span 類型是由使用的儀器庫自動確定的。OTel 定義了幾種不同的 Span 類型，每種類型對錯誤報告都有獨特的影響：

- **client**：用於發出的同步遠程呼叫（例如，發出的 HTTP 請求或 DB 呼叫）

- **server**：用於接收的同步遠程呼叫（例如，接收的 HTTP 請求或遠程過程呼叫）

- **internal**：用於不跨越process邊界的操作（例如，對函數呼叫進行監測）

- **producer**：用於創建可能稍後異步處理的作業（例如，插入隊列的工作）

- **consumer**：用於處理由生產者創建的工作，其可能在生產者 Span 結束後很長時間才開始



Span 類型是由使用的檢測工具自動判定的，雖然也能手動設定就是。



Span 還可以通過設置 Span 狀態進行增強。默認情況下，Span 狀態為 **UNSET**，除非另有指定。您可以將 Span 狀態標記為 **Error**，以指示 Span 描述了一個錯誤，或者將其標記為 **Ok**，以指示 Span 沒有錯誤。



### **Enhancing spans with span events**



另外，Span 還可以通過 Span 事件進行增強，Span 事件是嵌入到 Span 中的結構化日誌資訊。當 Span 狀態設置為 Error 時，將自動創建一個 Span 事件，該事件捕獲 Span 的錯誤訊息和追蹤堆疊。您可以通過向 Span 事件添加屬性來進一步增強這個錯誤。 OTel 提供了許多功能來處理和報告錯誤，使得開發人員能夠更加有效地監控和管理應用程序的錯誤情況。



這段內容對應到本書的 Ch 5-4 分享 Trace 有相關的介紹。



Jaeger 提供了有關查詢操作追蹤的一些操作以及可視化。下圖的紅點代表有 Span  被標記成 Error。

![img][圖3]



當我們點進去發生錯誤的Span時，就能看見如下圖的 log常見的 exception的資訊。

![img][圖4]


## Logs or spans to capture errors?



聽起來要用日誌還是Span來紀錄錯誤信號以及對應的事件，都可以做到！這很大程度取決於您的團隊以及現有遙測資料服務來決策。



Span 看起來對於捕獲錯誤很有效，因為將 Span 的狀態設定為 Error，很容易被篩選出來。但如果系統短時間內產生大量的 Span，而團隊又沒對 Span 進行 filter 或是尾部取樣，那麼我們會收到一堆 Span 來讓我們在排查問題時，只會比一堆日誌更加的困難。因為數量太多變成要找到真正的錯誤會變得困難，會變成一堆正常完成的 Span 操作，以及一堆非相關的 Span，使得要找到可能對系統的運行產生重大影響的錯誤非常困難。此外，一堆 Span 相比於日誌更佔用系統資源，包括儲存空間以及處理能力。當問題發生時，但資源卻不足，此時會沒辦法即時的排查問題。



使用 Span 作為捕捉錯誤的手段，還有個好處就是能搭配 Span Event，來紀錄異常訊息或其他你覺得相關的context內容。



但最重要的還是，我們的遙測資料後端是否同時支援日誌跟 Trace ? 以及兩者在查詢以及操作時是否很容易？兩者之間能否在遙測資料後端上實現相互關聯？



此外，不只是能將所有遙測資料轉發給OTel Collector 來進行 Filter的動作。如果團隊確定想更改OTel 檢測工具的默認錯誤處理行為，能通過 **[SetHandler](https://opentelemetry.io/docs/specs/otel/error-handling/#configuring-error-handlers) ** 來替換。 



最後，能根據團隊的具體需求以及遙測資料後端支持的功能，來選擇適合團隊現在的紀錄錯誤的方式。


原文參考自 [Dude, where's my error? How OpenTelemetry records errors](https://opentelemetry.io/blog/2024/otel-errors/)


[圖1]: ../images/otel_error_1.png
[圖2]: ../images/otel_error_2.png
[圖3]: ../images/otel_error_3.png
[圖4]: ../images/otel_error_4.png
