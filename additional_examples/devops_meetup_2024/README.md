# Starting with OpenTelemetry


## 2024 DevOps Mettup 

---

[投影片連結](https://docs.google.com/presentation/d/1eu_OP2N4e1vvLtWDiL4J2UKUuqhhXYsX/edit?usp=sharing&ouid=103239762195549851238&rtpof=true&sd=true)

[Demo 程式](https://github.com/tedmax100/devops_meetup)

## 使用方式
1. 編譯出 docker image
```make build```
2. 啟動系統
```docker compose -f compose.yml up -d```
3. 每次請務必關閉系統並刪除容器, 我 PostgreSQL 抄寫還沒弄好 QQ
```docker compose -f compose.yml down```
