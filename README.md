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

