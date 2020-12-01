gcloud config set project workkami
gcloud builds submit --tag gcr.io/workkami/go-mysql-api
gcloud beta run deploy --image gcr.io/workkami/go-mysql-api  --region asia-southeast1 --platform managed --allow-unauthenticated