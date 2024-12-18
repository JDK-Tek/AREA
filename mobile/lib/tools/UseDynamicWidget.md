# Documentation : Widget `Dynamic`

Le widget `Dynamic` est conçu pour afficher dynamiquement différents composants UI selon la valeur de son paramètre `title`. Il prend également un paramètre optionnel `extraParams` pour personnaliser davantage chaque composant.

---

## Constructeur

```dart
Dynamic({
  Key? key,
  required String title,
  Map<String, dynamic>? extraParams,
})
```

### **Paramètres**
- **`title`** (obligatoire) :  
  Détermine le type de widget à afficher. Voici les valeurs possibles :
  - `"textfield"`
  - `"timestamp"`
  - `"button"`
  - `"switch"`
  - `"dropdown"`
  - `"slider"`
  - `"date_picker"`
  - `"checkbox"`
  - `"listview"`
  - `"dialog"`
  
- **`extraParams`** (optionnel) :  
  Un `Map<String, dynamic>` contenant des paramètres spécifiques pour personnaliser le widget.

---

## Utilisation par `title`

### **1. TextField**
Affiche un champ de texte.  
**Exemple :**
```dart
Dynamic(
  title: "textfield",
  extraParams: {
    "labelText": "Entrez votre nom",
    "keyboardType": "text"
  },
)
```
- **`labelText`** : Texte du label (par défaut : "Entrez du texte").
- **`keyboardType`** : Type de clavier souhaité (par défaut : "text").

---

### **2. Timestamp**
Affiche un sélecteur de temps.  
**Exemple :**
```dart
Dynamic(
  title: "timestamp",
  extraParams: {
    "initialTime": TimeOfDay(hour: 10, minute: 30),
  },
)
```
- **`initialTime`** : Heure initiale (par défaut : `TimeOfDay.now()`).

---

### **3. Button**
Affiche un bouton avec une action.  
**Exemple :**
```dart
Dynamic(
  title: "button",
  extraParams: {
    "buttonText": "Soumettre",
    "onPressed": () {
      print("Bouton cliqué");
    },
  },
)
```
- **`buttonText`** : Texte du bouton (par défaut : "Cliquez Moi").
- **`onPressed`** : Action déclenchée au clic (par défaut : affiche un `SnackBar`).

---

### **4. Switch**
Affiche un interrupteur.  
**Exemple :**
```dart
Dynamic(
  title: "switch",
  extraParams: {
    "title": "Activer la fonctionnalité",
  },
)
```
- **`title`** : Texte affiché à côté du switch (par défaut : "Activer/Désactiver").

---

### **5. Dropdown**
Affiche un menu déroulant avec des options.  
**Exemple :**
```dart
Dynamic(
  title: "dropdown",
  extraParams: {
    "initialValue": "Option 2",
    "items": ["Option 1", "Option 2", "Option 3"],
  },
)
```
- **`initialValue`** : Valeur initiale sélectionnée (par défaut : "Option 1").
- **`items`** : Liste des options disponibles (par défaut : `["Option 1", "Option 2", "Option 3"]`).

---

### **6. Slider**
Affiche un curseur réglable.  
**Exemple :**
```dart
Dynamic(
  title: "slider",
  extraParams: {
    "initialValue": 0.7,
    "min": 0.0,
    "max": 1.0,
    "divisions": 20,
    "label": (value) => "${(value * 100).toStringAsFixed(0)}%",
  },
)
```
- **`initialValue`** : Valeur initiale (par défaut : `0.5`).
- **`min` / `max`** : Valeurs minimale et maximale (par défaut : `0.0` et `1.0`).
- **`divisions`** : Nombre de divisions (par défaut : `10`).
- **`label`** : Fonction pour formater l'étiquette affichée.

---

### **7. Date Picker**
Affiche un bouton qui ouvre un sélecteur de date.  
**Exemple :**
```dart
Dynamic(
  title: "date_picker",
  extraParams: {
    "buttonText": "Choisir une date",
    "initialDate": DateTime(2024, 1, 1),
    "firstDate": DateTime(2000),
    "lastDate": DateTime(2100),
  },
)
```
- **`buttonText`** : Texte du bouton (par défaut : "Choisir une date").
- **`initialDate` / `firstDate` / `lastDate`** : Plage des dates disponibles.

---

### **8. Checkbox**
Affiche une case à cocher.  
**Exemple :**
```dart
Dynamic(
  title: "checkbox",
  extraParams: {
    "title": "Accepter les termes",
    "initialValue": true,
  },
)
```
- **`title`** : Texte affiché à côté de la case (par défaut : "Activer l'option").
- **`initialValue`** : Valeur initiale (par défaut : `false`).

---

### **9. ListView**
Affiche une liste déroulante d’éléments.  
**Exemple :**
```dart
Dynamic(
  title: "listview",
  extraParams: {
    "items": ["Article 1", "Article 2", "Article 3"],
  },
)
```
- **`items`** : Liste des éléments (par défaut : `["Élément 0", ..., "Élément 9"]`).

---

### **10. Dialog**
Affiche un bouton qui ouvre une boîte de dialogue.  
**Exemple :**
```dart
Dynamic(
  title: "dialog",
  extraParams: {
    "title": "Confirmation",
    "content": "Voulez-vous continuer ?",
    "cancelText": "Non",
    "confirmText": "Oui",
    "onConfirm": () {
      print("Action confirmée");
    },
  },
)
```
- **`title`** : Titre de la boîte de dialogue.
- **`content`** : Contenu de la boîte.
- **`cancelText`** / **`confirmText`** : Texte des boutons.
- **`onConfirm`** : Action à exécuter sur confirmation.

---
