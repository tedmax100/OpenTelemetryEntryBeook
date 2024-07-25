-- script.lua
wrk.method = "GET"
wrk.headers["X-Scope-OrgID"] = "tenant1"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

wrk.path = "/loki/api/v1/query_range?query=sum%20by%20(status)%20(count_over_time({job=%22fluentbit%22}%20%7C%20json%20%7C%20line_format%20%22{{.log}}%20{{.container_name}}%22%20%7C%20pattern%20%60<ip>%20-%20-%20<_>%20%22<method>%20<uri>%20<_>%22%20<status>%20<size>%20<_>%20%22<agent>%22%20<_>%60%20[1m]))"

