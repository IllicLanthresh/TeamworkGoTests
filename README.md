The problem requires for this routine to be able to perform correctly on a small machine. 
Since we are using go I've chosen to take full advantage of the CPU power and goroutines and use as little memory as possible in the process

## CSV reading and email extraction
This approach uses a generator like pattern so the csv parsing and regex matching is done in a separate thread to prevent main thread blocking.
The generator checks for several possible failures before doing any work and then returns a channel that will spit out email addresses that are ready to use and validated

## Counter
The domain counter routine is really simple, it iterates over the generator channel and for each email address we extract the domain.

The domain gets stored in a map holding the count of each domain usage and if it's not present in the map, it gets fed to the sorter.
After getting the sorted list of domains, it's just mapping usage counts with each domain and returning that sorted list with the counts.

## Radix
This is where the rubber meets the road, my approach is an implementation of the concept behind Radix sort. 
It's programmed in such a way that it can be used by other modules, and it's not tied at all with the concept of email addresses nor domains.

A sorter gets feed the strings it must sort, those strings, as they're added, get placed in a bucket identified by the character at the cursor of that sorter.
Then, each bucket gets fed to a sorter, the cursor of that sorter gets advanced, and the cycle repeats.

To deal with shorter length strings when the cursor gets past their length, those strings get placed in a special bucket identified with a nullbyte. 
This is to take advantage of UTF-8 strings which don't have a character representation for code point zero(\x00)

After placing the longest strings in their respective buckets, going from the lowest level in recursion to the top, 
we iterate over the buckets in ascending order.
This order uses the byte value of the characters identifying the bucket(again taking advantage of UTF-8)

After all bucket levels get gathered, we end up with a sorted array of the initial data.

## Tests & Benchmarks
Multiple tests and benchmarks on the different layers of the module have been added, and they pass(I'm pretty sure you will have a lot more ideas for tests than I had).
What I was more interested in, was in the performance of the solution.

I've added benchmarks on the sorter itself creating randomized lists of different lengths, and with or without duplicates. 
I created those string lists using a really simple python oneliner to output the cartesian product of the alphabet characters.

Then I've added benchmarks to the importer+sorter layer using two sets of email addresses, yours and one I've created myself that is 1 m lines long.

I tried to find timing metrics of similar problems over the internet, but I couldn't find anything valuable. So my intuition is all I can say I used to evaluate this results.

### Sorter on its own
The results show what I was expecting, O(N logN), the more data you feed it, the more efficient it becomes, and the memory usage is almost negligible

### Importer: sorter with emails and counts
Since we are removing duplicates from the sorting algorithm, N gets downsized to log N, so I'm expecting the solution to become O(logN logN). 
In this case, after reviewing the results, I can confidently say that O(N) is a good approximation since it looks like the log ratio of N for each part gets the total really close to N.

Some metrics:
- 3k lines: 4 ms
- 1m lines: 150 ms

## Possible Improvements
### Regex
After reviewing the profiler data, I can see that the most CPU usage comes from the regex pattern matching(~60%), I could investigate some other solutions to replace the regex
### Non-blocking addition
I could refactor the sorter code, so it can get a stream of strings without blocking the main thread and process these strings in a goroutine
### Bubble sort
Bubble sort could also be an excellent candidate for this parallelized sorting pattern 
