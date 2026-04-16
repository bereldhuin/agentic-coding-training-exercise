# Plan d'Implémentation - Filtre Garantie pour Backend Python

## Contexte

Lebonpoint est une application de petites annonces multi-langages avec une architecture hexagonale (ports & adapters). Actuellement, le système de filtrage permet de filtrer les items par status, category, city, postal_code, is_featured, et delivery_available.

L'objectif est d'ajouter un nouveau filtre "garantie" (warranty) avec les valeurs possibles :
- Aucune garantie
- 6 mois
- 1 an
- 2 ans

**Décisions prises :**
- Implémentation **uniquement pour le backend Python** (pas les autres backends)
- Stockage en base : **colonne INTEGER stockant le nombre de mois** (0, 6, 12, 24)
- Le champ sera **optionnel (nullable)** pour permettre une migration progressive

## Architecture Actuelle - Flux de Filtrage

```
HTTP Request (FastAPI)
  ↓
Routes (query params extraction) → routes.py
  ↓
Use Case (validation & transformation) → list_items.py
  ↓
Repository Port (FilterOptions interface) → item_repository.py
  ↓
Repository Implementation (SQL WHERE clause) → sqlite_item_repository.py
  ↓
Database (SQLite)
```

## Fichiers Critiques à Modifier

### 1. Base de données
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/scripts/db-init.ts`
- Ajouter colonne `garantie_months INTEGER CHECK(garantie_months >= 0)` à la table `items`
- Ajouter index `idx_items_garantie_months` pour optimiser les requêtes filtrées
- **Note**: Pas de valeur DEFAULT car le champ est nullable
- Valeurs attendues : 0 (aucune garantie), 6, 12, 24

### 2. Domain Layer - Pas d'Enum nécessaire
**Note**: Pas besoin de créer un enum, on utilise directement `int | None` pour représenter le nombre de mois de garantie
- Valeurs valides : 0, 6, 12, 24 (validation au niveau de l'API avec Pydantic)

### 3. Domain Layer - Entity
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/python/src/domain/entities/item.py`
- Ajouter champ `garantie_months: int | None` à la classe `Item` (ligne ~29)
- Ajouter champ `garantie_months: int | None = None` à `CreateItemData` (ligne ~47)
- Ajouter champ `garantie_months: int | None = None` à `UpdateItemData` (ligne ~65)
- Ajouter champ `garantie_months: int | None = None` à `PatchItemData` (ligne ~83)
- Modifier `item_from_row()` pour récupérer la colonne `garantie_months` (ligne ~102)

### 4. Domain Layer - Repository Port
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/python/src/domain/repositories/item_repository.py`
- Ajouter `garantie_months: int | None = None` à la dataclass `FilterOptions`
- **Note**: Le filtre accepte un entier représentant le nombre de mois

### 5. Infrastructure Layer - HTTP Models
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/python/src/infrastructure/http/models.py`
- Ajouter `garantie_months: int | None = Field(None, ge=0)` à `ItemResponseModel` (ligne ~35)
- Ajouter `garantie_months: int | None = Field(None, ge=0)` à `CreateItemRequestModel` (ligne ~52)
- Ajouter `garantie_months: int | None = Field(None, ge=0)` à `UpdateItemRequestModel` (ligne ~69)
- Ajouter `garantie_months: int | None = Field(None, ge=0)` à `PatchItemRequestModel` (ligne ~86)
- Modifier `to_response_model()` pour inclure garantie_months (ligne ~138)
- Modifier `to_create_data()` pour mapper garantie_months (ligne ~185)
- Modifier `to_update_data()` pour mapper garantie_months (ligne ~240)
- Modifier `to_patch_data()` pour mapper garantie_months (ligne ~280)
- **Note**: `Field(None, ge=0)` valide que la valeur est >= 0 si fournie

### 6. Infrastructure Layer - Routes
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/python/src/infrastructure/http/routes.py`
- Ajouter query parameter `garantie_months: int | None = Query(None, ge=0, description="Filter by warranty period in months")` (ligne ~56)
- Passer garantie_months dans le dict `params` (ligne ~72)

### 7. Application Layer - Use Case
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/python/src/application/use_cases/list_items.py`
- Extraire `garantie_months` des query params (ligne ~50)
- Ajouter `garantie_months=garantie_months` dans la construction de `FilterOptions` (ligne ~50)

### 8. Infrastructure Layer - Repository Implementation
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/python/src/infrastructure/persistence/sqlite_item_repository.py`
- Ajouter condition WHERE pour garantie_months dans `find_all()` (après ligne ~133):
  ```python
  if filters.garantie_months is not None:
      if filters.garantie_months == 0:
          # Filtre 0 (aucune garantie) retourne aussi les valeurs NULL (non spécifiées)
          where_conditions.append("(garantie_months = ? OR garantie_months IS NULL)")
          params.append(filters.garantie_months)
      else:
          # Autres filtres: match exact uniquement
          where_conditions.append("garantie_months = ?")
          params.append(filters.garantie_months)
  ```
- Ajouter mapping garantie_months dans `create()` (ligne ~66): inclure garantie_months dans INSERT
- Ajouter mapping garantie_months dans `update()`: inclure garantie_months dans UPDATE
- Ajouter mapping garantie_months dans `patch()`: inclure garantie_months si présente dans UPDATE

### 9. Script de Seed
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/scripts/db-seed.ts`
- Ajouter array `GARANTIES_MONTHS = [0, 6, 12, 24]`
- Générer valeur aléatoire de garantie (en mois) pour chaque item de seed
- Inclure dans l'INSERT statement comme `garantie_months`

### 10. Frontend - HTML Form
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/client/index.html`
- Ajouter un champ `<select name="garantie_months">` dans le formulaire de filtres (après ligne ~307, avant la fermeture du `<div class="grid">`)
- Structure du select :
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
- **Note** : L'attribut `data-type="number"` permet à `readForm()` de convertir automatiquement la valeur en nombre (ligne ~830)

### 11. Frontend - i18n Français
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/client/i18n/fr.json`
- Ajouter dans `"filters"` :
  ```json
  "garantieLabel": "Garantie"
  ```
- Ajouter nouvelle section `"garantieOptions"` après `"categoryOptions"` (ligne ~49) :
  ```json
  "garantieOptions": {
    "none": "aucune garantie",
    "6_months": "6 mois",
    "1_year": "1 an",
    "2_years": "2 ans"
  }
  ```

### 12. Frontend - i18n Anglais
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/client/i18n/en.json`
- Ajouter dans `"filters"` :
  ```json
  "garantieLabel": "Warranty"
  ```
- Ajouter nouvelle section `"garantieOptions"` après `"categoryOptions"` (ligne ~49) :
  ```json
  "garantieOptions": {
    "none": "no warranty",
    "6_months": "6 months",
    "1_year": "1 year",
    "2_years": "2 years"
  }
  ```

### 13. Frontend - i18n Italien
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/client/i18n/it.json`
- Ajouter dans `"filters"` :
  ```json
  "garantieLabel": "Garanzia"
  ```
- Ajouter nouvelle section `"garantieOptions"` :
  ```json
  "garantieOptions": {
    "none": "nessuna garanzia",
    "6_months": "6 mesi",
    "1_year": "1 anno",
    "2_years": "2 anni"
  }
  ```

### 14. Frontend - i18n Occitan
**Fichier**: `/home/frederic-prost/tmp/agentic-coding-training-exercise/client/i18n/oc.json`
- Ajouter dans `"filters"` :
  ```json
  "garantieLabel": "Garantida"
  ```
- Ajouter nouvelle section `"garantieOptions"` :
  ```json
  "garantieOptions": {
    "none": "cap de garantida",
    "6_months": "6 meses",
    "1_year": "1 an",
    "2_years": "2 ans"
  }
  ```

## Ordre d'Implémentation Recommandé

### Backend Python
1. **Base de données** : Modifier `scripts/db-init.ts` pour ajouter la colonne `garantie_months INTEGER` et l'index
2. **Domain Entity** : Ajouter champ `garantie_months: int | None` aux dataclasses dans `item.py`
3. **Repository Port** : Ajouter `garantie_months` à `FilterOptions` dans `item_repository.py`
4. **HTTP Models** : Ajouter `garantie_months` aux Pydantic models dans `models.py` avec validation `ge=0`
5. **Routes** : Ajouter query parameter `garantie_months` dans `routes.py`
6. **Use Case** : Extraire et passer `garantie_months` dans `list_items.py`
7. **Repository Impl** : Implémenter filtrage SQL dans `sqlite_item_repository.py` (0 inclut NULL)
8. **Seed Script** : Ajouter données de garantie (0, 6, 12, 24) dans `db-seed.ts`

### Frontend
9. **HTML Form** : Ajouter select `garantie_months` (valeurs: 0, 6, 12, 24) dans `client/index.html`
10. **i18n FR** : Ajouter traductions françaises dans `client/i18n/fr.json`
11. **i18n EN** : Ajouter traductions anglaises dans `client/i18n/en.json`
12. **i18n IT** : Ajouter traductions italiennes dans `client/i18n/it.json`
13. **i18n OC** : Ajouter traductions occitanes dans `client/i18n/oc.json`

## Challenges et Points d'Attention

### Challenge 1 : Nullable vs Default Value
**Décision** : Le champ sera nullable (pas de DEFAULT dans la DB)
- Items existants auront `garantie = NULL`
- Le filtre ne retournera que les items avec une garantie explicitement définie
- Alternative rejetée : DEFAULT 'none' (aurait forcé tous les items existants à avoir "aucune garantie")

### Challenge 2 : Cohérence Multi-Backend
**Observation** : Cette implémentation ne sera que pour Python
- Les autres backends (TypeScript, Go, Kotlin, Swift) n'auront pas le filtre
- Le script de vérification cross-backend (`npm run verify`) échouera probablement
- **Justification** : Choix explicite de l'utilisateur pour accélérer le développement

### Challenge 3 : Migration des Données Existantes
**Impact** : Tous les items existants auront `garantie = NULL`
- **Solution implémentée** : Le filtre `garantie=none` retournera aussi les items avec `NULL` (via condition `OR garantie IS NULL`)
- Les autres filtres (`6_months`, `1_year`, `2_years`) ne retournent que les items avec valeur explicite
- Cela permet une migration progressive sans impact sur l'existant
- Les items anciens sont considérés comme "sans garantie" ce qui est logique du point de vue métier

### Challenge 4 : API OpenAPI Spec
**Note** : Le fichier `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/api/openapi.yaml` devrait être mis à jour mais n'a pas été inclus dans le périmètre initial
- **Recommandation** : Mettre à jour le spec OpenAPI pour documenter le nouveau filtre

### Challenge 5 : Frontend et i18n
**Solution implémentée** : Ajout complet du filtre dans le frontend
- Nouveau select dans le formulaire HTML (après delivery_available)
- Traductions dans les 4 langues supportées (FR, EN, IT, OC)
- Le système i18n existant (`data-i18n`) gère automatiquement l'affichage
- La fonction `readForm()` existante extrait automatiquement la valeur
- **Avantage** : Fonctionnalité complète de bout en bout, pas seulement API

## Plan de Vérification

### 1. Réinitialiser et Seed la Base
```bash
npm run db:init
npm run db:seed
npm run db:verify
```

### 2. Démarrer le Serveur Python
```bash
cd server/python
python3 -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python3 -m uvicorn src.main:app --reload --port 3001
```

### 3. Tests Manuels via API
```bash
# Liste tous les items
curl "http://localhost:3001/v1/items"

# Filtre par garantie = 0 mois (doit inclure les items avec NULL)
curl "http://localhost:3001/v1/items?garantie_months=0"

# Filtre par garantie = 12 mois (uniquement items avec garantie explicite de 1 an)
curl "http://localhost:3001/v1/items?garantie_months=12"

# Créer un item avec garantie de 24 mois
curl -X POST "http://localhost:3001/v1/items" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Item avec Garantie 2 ans",
    "price_cents": 50000,
    "condition": "new",
    "garantie_months": 24
  }'

# Créer un item sans garantie (NULL)
curl -X POST "http://localhost:3001/v1/items" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Item sans Garantie",
    "price_cents": 30000,
    "condition": "good"
  }'

# Créer un item avec garantie de 0 mois (explicitement sans garantie)
curl -X POST "http://localhost:3001/v1/items" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Item garantie explicite 0",
    "price_cents": 20000,
    "condition": "fair",
    "garantie_months": 0
  }'

# Combiner avec d'autres filtres
curl "http://localhost:3001/v1/items?garantie_months=12&status=active&category=Ordinateurs"
```

**Tests de validation** :
1. Vérifier que `garantie_months=0` retourne bien les items anciens (avec NULL) ET ceux avec 0 explicite
2. Vérifier que `garantie_months=12` ne retourne que les items avec exactement 12 mois
3. Créer un item sans garantie et vérifier qu'il apparaît avec le filtre `garantie_months=0`
4. Tester validation : essayer de créer un item avec `garantie_months=-1` (doit échouer)

### 4. Tests Unitaires
Vérifier que les tests existants passent toujours :
```bash
cd server/python
pytest tests/
```

**Note** : Des tests unitaires existants pourraient échouer si garantie devient obligatoire dans les fixtures

## Estimation de Complexité

- **Taille** : Moyenne (13 fichiers à modifier, ~220 lignes de code à ajouter/modifier)
  - Backend Python : 7 fichiers (pas d'enum à créer)
  - Frontend : 5 fichiers (1 HTML + 4 i18n JSON)
  - Scripts : 1 fichier (db-seed)
- **Risque** : Faible (pattern bien établi, stockage simple en INTEGER, changement isolé au backend Python + frontend)
- **Réversibilité** : Haute (simple rollback de migration DB + retrait du code)
- **Avantage INTEGER vs ENUM** : Plus flexible pour futures valeurs (ex: 18, 36 mois), pas de migration de schéma nécessaire

## Dépendances et Prérequis

- Base de données existante à `/home/frederic-prost/tmp/agentic-coding-training-exercise/server/database/db.sqlite`
- Python 3.x avec FastAPI, Pydantic, SQLite3
- Scripts Node.js (TypeScript) pour db-init et db-seed fonctionnels
