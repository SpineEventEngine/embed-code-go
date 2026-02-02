package org.example

// #docfragment "Main"
class Main private constructor() {
    // #enddocfragment "Main"
    companion object {
        // #docfragment "Main"
        fun main(args: Array<String>) {
            // #enddocfragment "Main"
            println("This is just a log message.")
            // #docfragment "Main"
            println(helperMethod())
        }
        // #enddocfragment "Main"
        private fun helperMethod(): String {
            return "42"
        }
    }
// #docfragment "Main"
}
// #enddocfragment "Main"
