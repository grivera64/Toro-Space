import React, { useEffect } from "react";

export default function Select() {
    const [users, setUsers] = React.useState([]);
    useEffect(() => {
        const fetchData = async () => {
            const response = await fetch('http://localhost:3030/account/self', {
                credentials: 'include'
            });
            const data = await response.json();
            setUsers(data);
        };
        fetchData();
    }, []);
    return (
        <div className='select container space-y-2 flex flex-col justify-center items-center mx-auto w-auto'>    
            <h1>Select a user</h1>
            <ul className='flex flex-row gap-2'>
                {
                    users.map((user, index) => (
                        <li className='box bg-gray-200 p-4 rounded-lg hover:bg-gray-400 cursor-pointer' key={index} onClick={
                            () => {
                                const selectUser = async () => {
                                    const resp = await fetch(`http://localhost:3030/account/self/user/${user.id}/select`, {
                                        method: 'PUT',
                                        credentials: 'include'
                                    });
                                    if (resp.status !== 200) {
                                        console.error('Failed to select user');
                                        return;
                                    }
                                    window.location.href = 'http://localhost:3000/home';
                                };
                                selectUser();
                            }
                        }>{user.display_name}</li>
                    ))
                }
            </ul>
        </div>
    );
}
