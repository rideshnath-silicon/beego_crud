<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>User List</title>
    <link rel="stylesheet" href="https://cdn.datatables.net/1.10.22/css/jquery.dataTables.min.css" />

    <!-- jQuery library file -->
    <script type="text/javascript" src="https://code.jquery.com/jquery-3.5.1.js">
    </script>

    <!-- Datatable plugin JS library file -->
    <script type="text/javascript" src="https://cdn.datatables.net/1.10.22/js/jquery.dataTables.min.js">
    </script>
    <style>
        table,
        th,
        td {
            border: 1px solid black;
            border-collapse: collapse;
        }
    </style>
</head>

<body>
    <h1>User List</h1>
    <div>
        <table id="users">
            <thead>
                <tr>
                    <th>SR</th>
                    <th>Name</th>
                    <th>Last Name</th>
                    <th>Email</th>
                    <th>Age</th>
                    <th>Country</th>
                    <th>Action</th>
                </tr>
            </thead>
            <tbody>
                {{range .users}}
                <tr> 
                    <td>{{.user_id}}</td>
                    <td>{{.name}}</td>
                    <td>{{.last_name}}</td>
                    <td>{{.user_email}}</td>
                    <td>{{.user_age}}</td>
                    <td>{{.country_name}}</td>
                    <td><a href=""></a></td>
                </tr>
                {{end}}
            </tbody>
        </table>
    </div>
</body>
<script>
    /* Initialization of datatable */
    $(document).ready(function () {
        $('#users').DataTable();
    }); 
</script>
</html>