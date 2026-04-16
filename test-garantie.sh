#!/bin/bash
# Script de test manuel pour la feature garantie
# Usage: ./test-garantie.sh

BASE_URL="http://localhost:8000"

echo "=========================================="
echo "Tests Feature Garantie - API Backend"
echo "=========================================="
echo ""

echo "1. Test: Tous les items (aucun filtre)"
echo "   GET $BASE_URL/v1/items"
curl -s "$BASE_URL/v1/items" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f'   → {len(data[\"items\"])} items retournés')"
echo ""

echo "2. Test: Items sans garantie (0 mois)"
echo "   GET $BASE_URL/v1/items?garantie_months=0"
curl -s "$BASE_URL/v1/items?garantie_months=0" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f'   → {len(data[\"items\"])} items'); [print(f'      - ID {item[\"id\"]}: {item[\"title\"]}') for item in data['items'][:3]]"
echo ""

echo "3. Test: Items avec 6 mois de garantie"
echo "   GET $BASE_URL/v1/items?garantie_months=6"
curl -s "$BASE_URL/v1/items?garantie_months=6" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f'   → {len(data[\"items\"])} items'); [print(f'      - ID {item[\"id\"]}: {item[\"title\"]}') for item in data['items'][:3]]"
echo ""

echo "4. Test: Items avec 1 an de garantie (12 mois)"
echo "   GET $BASE_URL/v1/items?garantie_months=12"
curl -s "$BASE_URL/v1/items?garantie_months=12" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f'   → {len(data[\"items\"])} items'); [print(f'      - ID {item[\"id\"]}: {item[\"title\"]}') for item in data['items'][:3]]"
echo ""

echo "5. Test: Items avec 2 ans de garantie (24 mois)"
echo "   GET $BASE_URL/v1/items?garantie_months=24"
curl -s "$BASE_URL/v1/items?garantie_months=24" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f'   → {len(data[\"items\"])} items'); [print(f'      - ID {item[\"id\"]}: {item[\"title\"]}') for item in data['items'][:3]]"
echo ""

echo "6. Test: Filtre combiné (garantie 6 mois + status active)"
echo "   GET $BASE_URL/v1/items?garantie_months=6&status=active"
curl -s "$BASE_URL/v1/items?garantie_months=6&status=active" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f'   → {len(data[\"items\"])} items'); [print(f'      - ID {item[\"id\"]}: {item[\"title\"]} (status: {item[\"status\"]})') for item in data['items'][:3]]"
echo ""

echo "7. Test: Validation - Valeur hors plage (999)"
echo "   GET $BASE_URL/v1/items?garantie_months=999"
curl -s "$BASE_URL/v1/items?garantie_months=999" | python3 -c "import sys, json; data = json.load(sys.stdin); print(f'   → {len(data[\"items\"])} items (attendu: 0)')"
echo ""

echo "8. Test: Validation - Type invalide (string)"
echo "   GET $BASE_URL/v1/items?garantie_months=abc"
HTTP_CODE=$(curl -s -w "%{http_code}" -o /tmp/garantie_test.json "$BASE_URL/v1/items?garantie_months=abc")
echo "   → HTTP $HTTP_CODE (attendu: 400)"
if [ "$HTTP_CODE" = "400" ]; then
  echo "   ✅ Validation correcte"
  cat /tmp/garantie_test.json | python3 -m json.tool 2>/dev/null | head -5
else
  echo "   ❌ Erreur: devrait retourner 400"
fi
echo ""

echo "=========================================="
echo "Distribution des garanties dans la DB"
echo "=========================================="
curl -s "$BASE_URL/v1/items" | python3 -c "
import sys, json
from collections import Counter
data = json.load(sys.stdin)
counts = Counter([item.get('garantie_months') for item in data['items']])
total = len(data['items'])
print(f'Total: {total} items')
for months in sorted(counts.keys()):
    count = counts[months]
    pct = (count / total * 100)
    print(f'  {months:2d} mois: {count:2d} items ({pct:5.1f}%)')
"
echo ""

echo "=========================================="
echo "Tests terminés!"
echo "=========================================="
