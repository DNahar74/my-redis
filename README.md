# Datatypes to be implemented

## Rules

- Data inside [] is required
- Data inside <> is optional


## Simple Strings

- Protocol Version -> RESP 2
- Category -> Simple
- First Byte -> +
- Format -> +[string(must not include \r or \n)]\r\n
- Leading and trailing whitespaces not allowed

## Simple Errors

- Protocol Version -> RESP 2
- Category -> Simple
- First Byte -> -
- Format -> -[simple string]\r\n
- Redis convention -> -<ERROR_NAME(all caps)>[simple string]\r\n
- Leading and trailing whitespaces not allowed

## Integers

- Protocol Version -> RESP2
- Category -> Simple
- First Byte -> :
- Format -> :<+|->[integer]\r\n

## Bulk Strings

- Protocol Version -> RESP2
- Category -> Aggregate
- First Byte -> $
- Format -> $[length]\r\n[value]\r\n
- Leading and trailing whitespaces are allowed

## Arrays

- Protocol Version -> RESP2
- Category -> Aggregate
- First Byte -> *
- Format -> *[number-of-elements]\r\n[element-1]...[element-n]
- Remember, no need to end it with a \r\n for the array (element already ends with \r\n)

## Nulls

- Protocol Version -> 
- Category -> 
- First Byte -> 
- Format -> 

## Booleans
- Protocol Version -> 
- Category -> 
- First Byte -> 
- Format -> 
