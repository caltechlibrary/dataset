SELECT '"' || _Key || '"'
FROM friends
WHERE src->>'given' LIKE 'Mojo' 
   OR src->>'family' LIKE 'Mojo'
