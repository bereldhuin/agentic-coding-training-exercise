# Diagnostic du problème d'affichage des traductions

## Problème rapporté
Le texte littéral `filters.garantieLabel` s'affiche sur la page au lieu de la traduction "Warranty" (ou "Garantie" en français).

## Cause racine identifiée

Le dictionnaire de fallback `FALLBACK_TRANSLATIONS` intégré dans `/client/index.html` était **incomplet**. Les clés suivantes manquaient pour les 4 langues (en, fr, it, oc):

1. `filters.garantieLabel` - Le label "Warranty"/"Garantie"
2. `garantieOptions.none` - "no warranty"/"aucune garantie"
3. `garantieOptions.6_months` - "6 months"/"6 mois"
4. `garantieOptions.1_year` - "1 year"/"1 an"
5. `garantieOptions.2_years` - "2 years"/"2 ans"

## Pourquoi le fallback était utilisé?

Le système i18n fonctionne comme suit:

1. Au chargement de la page, la fonction `loadTranslations(lang)` tente de charger `/i18n/{lang}.json` via fetch
2. Si le fetch réussit (`response.ok === true`), les traductions du fichier JSON sont utilisées
3. Si le fetch échoue (erreur réseau, 404, etc.) ou si `response.ok === false`, le fallback embarqué est utilisé

Les fichiers JSON (`/client/i18n/en.json`, etc.) contenaient bien toutes les clés de traduction, mais:
- Le fallback était incomplet
- En cas d'erreur de chargement des JSON, le fallback prenait le relais avec des traductions manquantes
- La fonction `t(key)` retourne la clé elle-même si la traduction n'existe pas (ligne 719 dans index.html)

## Vérifications effectuées

✓ Les fichiers JSON `/client/i18n/*.json` contiennent bien les clés `filters.garantieLabel` et `garantieOptions.*`
✓ Le serveur HTTP (`python3 -m http.server 8080`) sert bien ces fichiers avec status 200
✓ Le Content-Type est correct (`application/json`)

## Solution appliquée

Ajout des clés manquantes dans `FALLBACK_TRANSLATIONS` pour les 4 langues:

### Anglais (en)
```javascript
filters: {
  // ...
  garantieLabel: "Warranty"
},
garantieOptions: {
  none: "no warranty",
  "6_months": "6 months",
  "1_year": "1 year",
  "2_years": "2 years"
}
```

### Français (fr)
```javascript
filters: {
  // ...
  garantieLabel: "Garantie"
},
garantieOptions: {
  none: "aucune garantie",
  "6_months": "6 mois",
  "1_year": "1 an",
  "2_years": "2 ans"
}
```

### Italien (it)
```javascript
filters: {
  // ...
  garantieLabel: "Garanzia"
},
garantieOptions: {
  none: "nessuna garanzia",
  "6_months": "6 mesi",
  "1_year": "1 anno",
  "2_years": "2 anni"
}
```

### Occitan (oc)
```javascript
filters: {
  // ...
  garantieLabel: "Garantida"
},
garantieOptions: {
  none: "cap de garantida",
  "6_months": "6 meses",
  "1_year": "1 an",
  "2_years": "2 ans"
}
```

## Test de validation

Pour valider la correction:

1. Ouvrir http://localhost:8080 dans un navigateur
2. Vérifier que le label "Warranty" s'affiche au lieu de "filters.garantieLabel"
3. Vérifier que les options du dropdown affichent "Any", "no warranty", "6 months", "1 year", "2 years"
4. Changer la langue vers "Français" et vérifier que "Garantie" s'affiche
5. Vérifier les options: "Tous", "aucune garantie", "6 mois", "1 an", "2 ans"
6. Tester également en Italien et Occitan

## Fichiers modifiés

- `/home/frederic-prost/tmp/agentic-coding-training-exercise/client/index.html` - Ajout des clés manquantes dans FALLBACK_TRANSLATIONS

## Note technique

Le système de fallback est une bonne pratique qui permet de:
- Fonctionner même si les fichiers JSON externes ne se chargent pas
- Avoir une expérience dégradée mais fonctionnelle en cas de problème réseau
- Garantir que l'application affiche du contenu lisible même en cas d'échec

**Important**: À l'avenir, lors de l'ajout de nouvelles clés de traduction, il faut penser à mettre à jour:
1. Les 4 fichiers JSON dans `/client/i18n/*.json`
2. Les 4 dictionnaires de fallback dans `FALLBACK_TRANSLATIONS` dans `/client/index.html`
