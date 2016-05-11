# handlersocket
Go driver for MariaDB/MySQL handlersocket plugin

### Opening/closing a connection

        if hs, err = handlersocket.Connect("127.0.0.1:9999", ""); err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Successfully connected to the server")
        }

        defer hs.Close()

### Opening an index

        if idx, err = hs.OpenIndex(1, "PRIMARY", "test", "t1", []string{"i"}); err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("open_index operation successful, ", idx)
        }

### Find/Select operation

        if rows, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"1", true}}).
                Select(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Result set : ", rows)
        }

### Insert operation

        if err = idx.
                Reset().
                Insert([]handlersocket.NullString{{"3", true}}); err != nil {
                fmt.Println(err)
                os.Exit(1)
        }

### Update operation

        if res, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"3", true}}).
                Update([]handlersocket.NullString{{"2", true}}); err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Result : ", res)
        }

### Delete operation

        if res, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"2", true}}).
                Delete(); err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Result : ", res)
        }

### Increment operation

        res, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"2", true}}).
                Increment([]handlersocket.NullString{{"1", true}})
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Increment result : ", res)
        }

### Decrement operation

        res, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"3", true}}).
                Decrement([]handlersocket.NullString{{"1", true}})
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Decrement result : ", res)
        }

### Updated operation

        rows, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"1", true}}).
                Updated([]handlersocket.NullString{{"2", true}})
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Updated records : ", rows)
        }

### Deleted operation

        rows, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"2", true}}).
                Deleted()
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Deleted records : ", rows)
        }

### Incremented operation

        rows, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"2", true}}).
                Incremented([]handlersocket.NullString{{"1", true}})
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Incremented records : ", rows)
        }

### Decremented operation

        rows, err = idx.
                Reset().
                Operator("=").
                Values([]handlersocket.NullString{{"3", true}}).
                Decremented([]handlersocket.NullString{{"1", true}})
        if err != nil {
                fmt.Println(err)
                os.Exit(1)
        } else {
                fmt.Println("Decremented records : ", rows)
        }

