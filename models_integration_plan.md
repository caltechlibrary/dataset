# Models Integration Plan for dataset

*Project: Add schema-based validation to dataset collections using the models module*
*Date: 2025-01-11*
*Status: Plan Approved, Implementation Pending*

---

## 🎯 Goal

Add schema-based validation to dataset collections using the models module, enabling structured data validation for CrossRef, DataCite, and other archival record formats. This implements the vision from `structures/README.md` for a turnkey backend system that validates records without custom middleware.

---

## ✅ Prerequisites (Already Complete)

| Task | Status | Location | Notes |
|------|--------|----------|-------|
| Models package identifier types | ✅ Done | `models/types.go` | ISBN, ISSN, DOI, ArXiv, ORCID, ROR, ISNI, UUID, PMID, PMCID, FundRef, LCNAF, VIAF, SNAC |
| Models package HTML5 types | ✅ Done | `models/types.go` | date, email, tel, url, text, textarea, checkbox, radio, number, range, color, etc. |
| Dataset Collection.Model field | ✅ Done | `dataset/collection.go` | `Model *models.Model` field exists |
| Dataset model.yaml loading | ✅ Done | `dataset/collection.go:155-176` | Loads `model.yaml` from collection root directory |
| Existing validation in API | ✅ Partial | `dataset/api_routes.go` | Uses Model.ValidateMapInterface for simple validation |

---

## 📌 Phase 1: Extend models for Nested Structures *(Foundation)*

**Objective:** Add support for nested objects and lists to the models package to handle complex records like CrossRef and DataCite.

### Tasks

| # | Task | File | Lines Est. | Status |
|---|------|------|------------|--------|
| 1.1 | Add `Elements []*Element` field to `Element` struct for nested elements | `models/element.go` | 5 | ✅ |
| 1.2 | Add `IsList bool` field to `Element` struct for array/list types | `models/element.go` | 3 | ✅ |
| 1.3 | Add `IsObject bool` field to `Element` struct for object/map types | `models/element.go` | 3 | ✅ |
| 1.4 | Update validation to recursively handle nested structures | `models/model.go` | 150 | ✅ |
| 1.5 | Update YAML parsing to populate nested `Elements` | `models/interactive.go` | 0 | ❌ (auto-works) |
| 1.6 | Add `GetNestedElement(path string) (*Element, bool)` helper | `models/model.go` | 25 | ✅ |
| 1.7 | Add `isRequired(elem *Element) bool` helper | `models/model.go` | 12 | ✅ |
| 1.8 | Update `Check()` to validate nested structures | `models/element.go` | 15 | ✅ |
| 1.9 | Add tests for nested validation | `models/model_test.go` | 80 | ✅ |

### Checkpoint 1
> ✅ **COMPLETE** - Models package supports hierarchical schemas matching CrossRef/DataCite structures.
> Nested `Element.Elements` allows defining objects within objects and lists of objects.

---

## 📌 Phase 2: Extend dataset Collection Configuration *(Configuration)*

**Objective:** Allow schemas to be defined in datasetd settings and referenced by collections.

### Tasks

| # | Task | File | Lines Est. | Status |
|---|------|------|------------|--------|
| 2.1 | Add `Schemas map[string]*models.Model` to Config struct | `dataset/config.go` | 5 | ✅ |
| 2.2 | Add `SchemaName string` field to Collection struct | `dataset/config.go` | 3 | ✅ |
| 2.3 | Add `Validate bool` field to Collection struct | `dataset/config.go` | 3 | ✅ |
| 2.4 | Add `ResolveSchemas()` method to resolve schema references | `dataset/config.go` | 30 | ✅ |
| 2.5 | Update `ConfigOpen` to call `ResolveSchemas()` | `dataset/config.go` | 5 | ✅ |
| 2.6 | Update `API.Init` to apply config schema to collection | `dataset/api.go` | 5 | ✅ |
| 2.7 | Maintain backward compatibility with existing `model.yaml` loading | `dataset/collection.go` | 0 | ✅ (auto-works) |

### Checkpoint 2
> ✅ **COMPLETE** - Settings.yaml can define global schemas. Collections can reference schemas by name or define inline schemas. Schema resolution happens automatically during config loading.

---

## 📌 Phase 3: Update Validation Logic *(Integration)*

**Objective:** Integrate schema validation into dataset API handlers.

### Tasks

| # | Task | File | Lines Est. | Status |
|---|------|------|------------|--------|
| 3.1 | Create `ValidateRecord(model *Model, data map[string]interface{}, path string) error` | `dataset/validation.go` (new) | 60 | ⬜ |
| 3.2 | Create `ValidationError` type with field path and message | `dataset/validation.go` | 20 | ⬜ |
| 3.3 | Update `createObjectHandler` to validate against collection schema | `dataset/api_routes.go` | 15 | ⬜ |
| 3.4 | Update `updateObjectHandler` to validate against collection schema | `dataset/api_routes.go` | 15 | ⬜ |
| 3.5 | Add validation error formatting with field paths | `dataset/validation.go` | 30 | ⬜ |
| 3.6 | Add `X-Validation-Errors` header for API error responses | `dataset/api_routes.go` | 10 | ⬜ |
| 3.7 | Handle empty/optional fields in validation | `dataset/validation.go` | 20 | ⬜ |

### Checkpoint 3
> API endpoints validate incoming data against collection schemas. Validation errors are returned with field paths.

---

## 📅 Quick Start (Minimal Viable Implementation)

For rapid delivery of core functionality, focus on these tasks in order:

### Priority Order

1. **Phase 1:** Tasks 1.1-1.3 (Add nested structure fields to Element) - *~10 lines*
2. **Phase 1:** Task 1.4 (Update ValidateMapInterface for recursion) - *~50 lines*
3. **Phase 1:** Task 1.9 (Basic nested validation tests) - *~40 lines*
4. **Phase 2:** Tasks 2.1-2.3 (Add schema config fields) - *~10 lines*
5. **Phase 2:** Task 2.4 (Parse global schemas from YAML) - *~40 lines*
6. **Phase 3:** Tasks 3.1-3.2 (Create validation.go) - *~80 lines*
7. **Phase 3:** Tasks 3.3-3.4 (Integrate into API handlers) - *~30 lines*

**Total Quick Start:** ~260 lines of code

**Result:** Working schema validation for nested structures via collection configuration.

---

## 📁 File Changes Summary

| File | Changes | Status |
|------|---------|--------|
| `models/element.go` | Add 3 fields for nesting, update Check() | ✅ |
| `models/model.go` | Add ValidateInterface(), ValidateRecursive(), validateList(), validateObject(), GetNestedElement(), hasNestedElement(), isRequired() | ✅ |
| `models/interactive.go` | YAML parsing for nested elements (auto-works with existing code) | ✅ |
| `models/model_test.go` | Add TestNestedValidation, TestGetNestedElement, TestElementCheck | ✅ |
| `dataset/config.go` | Add Schemas map, SchemaName, Validate fields, ResolveSchemas(), update ConfigOpen | ✅ |
| `dataset/api.go` | Apply config schema to collection after Open() | ✅ |
| `dataset/config.go` | Add Schemas map, schema resolution | ⬜ |
| `dataset/collection.go` | Add SchemaName, Validate fields | ⬜ |
| `dataset/validation.go` | New file: validation logic | ⬜ |
| `dataset/api_routes.go` | Integrate validation | ⬜ |
| `dataset/models_integration_plan.md` | This document | ✅ |

---

## 🧪 Testing Strategy

### Unit Tests
- Models package: nested structure parsing and validation
- Config parsing: schema loading and resolution
- Validation: all identifier types, nested paths, error messages

### Integration Tests
- API endpoints with schema-enabled collections
- Backward compatibility with existing model.yaml files
- Mixed collections (some with schemas, some without)

### Test Data
- CrossRef record examples
- DataCite record examples
- Flat records (existing behavior)
- Nested records (new behavior)

---

## 📚 Documentation Updates

| File | Update |
|------|--------|
| `dataset/datasetd.1.md` | Add schema configuration section |
| `dataset/README.md` | Add schema usage examples |
| `dataset/demo/settings.yaml` | Example with schemas |
| `dataset/demo/crossref_schema.yaml` | CrossRef schema example |

---

## 🔄 Backward Compatibility

All changes must maintain backward compatibility:

1. **Existing model.yaml files** continue to work (loaded from collection directory)
2. **Collections without schemas** continue to work (validation disabled)
3. **Existing API behavior** unchanged when no schema is configured
4. **Validation is opt-in** via `validate: true` in collection config

---

## 🎯 Success Criteria

- [ ] CrossRef record can be defined as a schema in YAML
- [ ] Schema can be referenced by multiple collections
- [ ] Validation rejects invalid CrossRef records
- [ ] Validation accepts valid CrossRef records
- [ ] Nested author lists validate correctly
- [ ] All identifier types (ISBN, ISSN, DOI, etc.) validate in nested structures
- [ ] Existing collections without schemas continue to work

---

## 📞 Notes

- The models package already has extensive identifier type support (ISBN, ISSN, DOI, ORCID, ROR, ISNI, UUID, PMID, PMCID, FundRef, LCNAF, VIAF, SNAC, plus all HTML5 types).
- The dataset package already loads model.yaml from collection directories.
- The API routes already have conditional validation using Model.ValidateMapInterface.
- The key missing piece is nested structure support in the models package.
- Once nested structures are supported, the rest is configuration and integration.

---

## 📌 Related Documents

- `structures/README.md` - Project vision and requirements
- `models/README.md` - Models package documentation
- `metadatatools/` - Identifier type reference implementation (TypeScript)
- `dataset/datasetd.1.md` - Current datasetd documentation
