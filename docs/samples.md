
# Samples

Two commands, _dataset keys_ and _dsfind_ support the `-sample` option. The sample option
expects a sample size as an argument. If the sample size is greater than zero then a sample
output will be taken.  On _dataset keys_ the sample is taken after any filter is supplied.
For _dsfind_ the sample is taken after the query results have been calculated. Is 
the sample size is greater then the results returned then the who results set is return
without random sampling.  If sample size is less than result set then a random sampling of
the results is taken.

If you are doing Machine Leaning type of sampling (e.g. calculating a test and training set)
then normally you create a *test* key list like this `dataset -sample="$N" keys` where `$N` 
holds the test sample size. After keylist is generated you can then create a training set
by excluding the keys associated with the sample.

