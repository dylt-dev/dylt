Login to Cloudflare 

https://dash.cloudflare.com/login

Create an API key

https://dash.cloudflare.com/profile/managed-profile/preferences



Test the API token
```
api_token=A2T6PYqkgqydQav4ftZUk9frOukk-CYZ1AU8Suxt
domain=dylt.dev


# Get zone id for domain
curl --silent https://api.cloudflare.com/client/v4/zones/ --header "Authorization: Bearer $api_token" | jq -r ".result[] | select(.name == \"$domain\") | .id"

# Get records for zone
curl --silent --header "Authorization: Bearer $api_token" https://api.cloudflare.com/client/v4/zones/$zone_id/dns_records | jq -r '.result[] | [.name, .type, .content, .comment] | @csv'

# Export zone file for a domain
curl --silent --header "Authorization: Bearer $api_token" https://api.cloudflare.com/client/v4/zones/$zone_id/dns_records/export
```

Other ...
```
# Get record id (eg so you can delete the record)
curl --silent --header "Authorization: Bearer $api_token" https://api.cloudflare.com/client/v4/zones/$zone_id/dns_records | jq -r '.result[] | select(.type=="A" and .name=="vm.dylt.dev") | .id'
```
More advanced usage ...

Visit the API page
https://dash.cloudflare.com/login
