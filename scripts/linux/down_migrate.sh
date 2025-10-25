
SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
cd "$SCRIPT_DIR/../.."
source .env

cd "$SCRIPT_DIR/../../db/schema"
pwd
goose postgres $MIGRATION_URL down
