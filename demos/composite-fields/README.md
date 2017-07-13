
# Hybrid Indexes

This demo explores building an index of actual content fields as well as creating
new fields from the record by applying a template to yielda new value.  In our 
example we are transforming the ORCID name structure into two indexed values
[family_name](family_name.tmpl)  and [given_name](given_name.tmpl). The final
index should be based on a final field called _orcid_ which corresponds to the
ORCID id in our harvested data.


