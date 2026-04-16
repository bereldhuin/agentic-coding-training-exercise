# Documentation Interface - Filtre Garantie

## Vue d'ensemble

Le filtre garantie a ete integre dans la section "Search & Filter" du frontend.

---

## Structure HTML du Filtre

```html
<label>
  <span data-i18n="filters.garantieLabel">Warranty</span>
  <select name="garantie_months" data-type="number">
    <option value="" data-i18n="common.any">Any</option>
    <option value="0" data-i18n="garantieOptions.none">No warranty</option>
    <option value="6" data-i18n="garantieOptions.6_months">6 months</option>
    <option value="12" data-i18n="garantieOptions.1_year">1 year</option>
    <option value="24" data-i18n="garantieOptions.2_years">2 years</option>
  </select>
</label>
```

---

## Position dans l'Interface

Le filtre garantie est positionne dans la grille de filtres, apres le filtre "Delivery available".

**Ordre des filtres**:
1. Search query
2. Status
3. Category
4. Min price (€)
5. Max price (€)
6. Postal code
7. Delivery available
8. **Warranty** ← NOUVEAU

---

## Traductions

### Anglais (EN)
```
Label: "Warranty"
Options:
  - "Any"
  - "No warranty"
  - "6 months"
  - "1 year"
  - "2 years"
```

### Francais (FR)
```
Label: "Garantie"
Options:
  - "Tous"
  - "Aucune garantie"
  - "6 mois"
  - "1 an"
  - "2 ans"
```

### Italien (IT)
```
Label: "Garanzia"
Options:
  - "Qualsiasi"
  - "Nessuna garanzia"
  - "6 mesi"
  - "1 anno"
  - "2 anni"
```

### Occitan (OC)
```
Label: "Garantia"
Options:
  - "Totes"
  - "Cap de garantia"
  - "6 meses"
  - "1 an"
  - "2 ans"
```

---

## Comportement Utilisateur

### Scenario 1: Filtre par garantie seule
1. L'utilisateur selectionne "6 months" dans le filtre Warranty
2. L'utilisateur clique sur "Search items"
3. L'API est appelee: `GET /v1/items?garantie_months=6`
4. 6 items sont affiches dans la section Results

### Scenario 2: Filtre combine
1. L'utilisateur selectionne "1 year" dans Warranty
2. L'utilisateur selectionne "active" dans Status
3. L'utilisateur clique sur "Search items"
4. L'API est appelee: `GET /v1/items?garantie_months=12&status=active`
5. Seuls les items actifs avec 1 an de garantie sont affiches

### Scenario 3: Changement de langue
1. L'utilisateur selectionne "Francais" dans le selecteur de langue
2. Le label "Warranty" devient "Garantie"
3. Les options sont traduites:
   - "No warranty" → "Aucune garantie"
   - "6 months" → "6 mois"
   - "1 year" → "1 an"
   - "2 years" → "2 ans"

---

## Integration avec le Systeme i18n

Le filtre utilise le systeme i18n existant avec:
- `data-i18n` attributes sur les elements HTML
- Fichiers JSON de traduction dans `/client/i18n/`
- Mise a jour dynamique lors du changement de langue

**Cles i18n utilisees**:
- `filters.garantieLabel` - Label du filtre
- `garantieOptions.none` - Option "Aucune garantie"
- `garantieOptions.6_months` - Option "6 mois"
- `garantieOptions.1_year` - Option "1 an"
- `garantieOptions.2_years` - Option "2 ans"
- `common.any` - Option "Tous/Any" (reutilisee)

---

## Traitement Frontend → Backend

### Frontend (JavaScript)
```javascript
// Lecture du formulaire
const data = readForm(listForm);
// data.garantie_months = "12" (string depuis le select)

// Conversion en nombre (data-type="number")
if (field.dataset.type === "number") {
  value = raw === "" ? undefined : Number(raw);
}
// garantie_months devient 12 (number)

// Construction de l'URL
const params = new URLSearchParams();
params.set("garantie_months", "12");
// GET /v1/items?garantie_months=12
```

### Backend (Python FastAPI)
```python
@router.get("/items")
async def list_items(
    garantie_months: Optional[int] = None,
    # ... autres parametres
):
    items = await repository.list_items(
        garantie_months=garantie_months
    )
    return {"items": items}
```

---

## Validation des Donnees

### Cote Frontend
- Type HTML: `<select>` avec valeurs predefinies
- data-type="number" pour conversion automatique
- Valeurs possibles: "", "0", "6", "12", "24"

### Cote Backend
- Type: `Optional[int]`
- Validation Pydantic automatique
- Erreur 400 si type invalide

**Exemple d'erreur**:
```
GET /v1/items?garantie_months=abc
→ 400 Bad Request
{
  "error": {
    "code": "validation_error",
    "message": "Validation failed",
    "details": {
      "query.garantie_months": "Input should be a valid integer"
    }
  }
}
```

---

## Accessibilite

- Label explicite avec `<span>` pour lecteurs d'ecran
- Element `<select>` natif pour navigation clavier
- Options traduites selon la langue de l'utilisateur
- Valeurs semantiques (0, 6, 12, 24 = mois)

---

## Responsive Design

Le filtre utilise la grille CSS existante:
```css
.grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 10px;
}
```

Sur mobile, le filtre s'empile verticalement avec les autres filtres.

---

## Tests Manuels Recommandes

### Test 1: Affichage initial
- [ ] Ouvrir http://localhost:8080
- [ ] Verifier la presence du filtre "Warranty"
- [ ] Verifier les 5 options (Any, No warranty, 6 months, 1 year, 2 years)

### Test 2: Filtrage fonctionnel
- [ ] Selectionner "6 months"
- [ ] Cliquer "Search items"
- [ ] Verifier l'affichage de 6 items

### Test 3: Changement de langue
- [ ] Selectionner "Francais"
- [ ] Verifier "Warranty" → "Garantie"
- [ ] Verifier "6 months" → "6 mois"

### Test 4: Filtre combine
- [ ] Selectionner "1 year" + Status "active"
- [ ] Cliquer "Search items"
- [ ] Verifier les resultats (items actifs avec 12 mois)

### Test 5: Reset
- [ ] Selectionner "Any" dans Warranty
- [ ] Cliquer "Search items"
- [ ] Verifier l'affichage de tous les items (15)

---

## URL du Frontend

**Fichier local**: `file:///home/frederic-prost/tmp/agentic-coding-training-exercise/client/index.html`

**Serveur HTTP** (pour tests):
```bash
cd /home/frederic-prost/tmp/agentic-coding-training-exercise/client
python3 -m http.server 8080
# Puis ouvrir: http://localhost:8080
```

**Serveur backend**: http://localhost:8000
