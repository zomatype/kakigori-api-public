#!/usr/bin/env bash
# かき氷APIを Google Cloud Run にデプロイする。
#
# 使い方:
#   gcloud auth login            # 先に認証しておく
#   ./deploy-cloudrun.sh <GCPプロジェクトID>
#
# 発表が終わったら以下で課金を止められる:
#   gcloud run services update kakigori-api --region asia-northeast1 --min-instances=0
#   gcloud run services delete kakigori-api --region asia-northeast1
set -euo pipefail

PROJECT_ID="${1:-$(gcloud config get-value project 2>/dev/null)}"
REGION="${REGION:-asia-northeast1}"
SERVICE="${SERVICE:-kakigori-api}"

if [[ -z "$PROJECT_ID" || "$PROJECT_ID" == "(unset)" ]]; then
  echo "エラー: GCPプロジェクトIDを指定してください: ./deploy-cloudrun.sh <PROJECT_ID>" >&2
  exit 1
fi

echo "==> プロジェクト: $PROJECT_ID / リージョン: $REGION / サービス: $SERVICE"
gcloud config set project "$PROJECT_ID" >/dev/null

echo "==> 必要なAPIを有効化"
gcloud services enable run.googleapis.com cloudbuild.googleapis.com artifactregistry.googleapis.com

# --max-instances=1 は必須。このAPIは注文をインメモリに持つので、
# 複数インスタンスに分かれると「注文したのにレシート画面で404」が発生する。
# --min-instances=1 でインスタンスを常駐させ、コールドスタートと状態消失を防ぐ。
echo "==> デプロイ中（初回はCloud Buildで数分かかります）"
gcloud run deploy "$SERVICE" \
  --source . \
  --region "$REGION" \
  --platform managed \
  --allow-unauthenticated \
  --min-instances=1 \
  --max-instances=1 \
  --concurrency=80 \
  --cpu=1 \
  --memory=512Mi \
  --timeout=30s \
  --set-env-vars="^@^KAKIGORI_STORE_IDS=store-001,store-002@KAKIGORI_MAX_ORDERS=1000@KAKIGORI_STORE_MAX_REQUESTS=100000"

URL="$(gcloud run services describe "$SERVICE" --region "$REGION" --format='value(status.url)')"

echo
echo "==> 動作確認"
curl -fsS "$URL/health" && echo
curl -fsS "$URL/v1/stores/store-001/menu" >/dev/null && echo "メニュー取得 OK"

echo
echo "============================================================"
echo " デプロイ完了"
echo "   API URL: $URL"
echo
echo " 次の手順: Vercel の環境変数を更新して再デプロイ"
echo "   VITE_API_BASE = $URL"
echo "============================================================"
