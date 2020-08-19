# Register your app on the azure portal
# https://portal.azure.com/#blade/Microsoft_AAD_RegisteredApps/ApplicationsListBlade
# and extract the information once the registration is complete
export AZURE_CLIENT_ID=""
export AZURE_CLIENT_SECRET=""
export AZURE_TENANT_ID=""
export WEB_URL="http://localhost:8080"
export AZURE_REDIRECT_URL="/auth/azure/redirect"
export AZURE_DOMAIN=""  # The domain who can access the app like "outlook.com" or "company.com"
export AZURE_RESOURCE="00000003-0000-0000-c000-000000000000"

# Session specific setting
export SESSION_NAME="TEST-SESSION"   # Name your session
export JWT_NAME="TEST-JWT"       # Name your JWT
export JWT_TOKEN="TEST-HIGHLY-SECURE-JWT-TOKEN"      # Provide a secret which which will secure the JWT
export EXPIRY_DATE="12"  # The expiration in hours after which the JWT token needs to be renewed

# Generic
export API_PORT="8080"
export LOG_LEVEL="DEBUG"

# Run the app
go run *.go