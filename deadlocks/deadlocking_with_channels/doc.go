package main

//It’s important to understand that deadlocks aren’t limited to the use of mutexes.
//Deadlocks can occur whenever executions hold mutually exclusive resources and
//request other ones—this also applies to channels.
