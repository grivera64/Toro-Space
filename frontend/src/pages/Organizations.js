import React from "react";

import FilterSearch from "../components/common/FilterSearch";

export default function Organizations() {
    const [organizations, setOrganizations] = React.useState([]);
    const [searchQuery, setSearchQuery] = React.useState('');

    React.useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:3030/organizations?search_query=' + searchQuery);
                const data = await response.json();
                setOrganizations(data);
            } catch (error) {
                setOrganizations([]);
            }
        };
        fetchData();
    }, [searchQuery]);

    return (
        <div className='topics-page w-full h-auto'>
            <div className='justify-center m-auto flex flex-col gap-3 mt-3'>
                <p className='text-5xl text-center'>Organizations</p>
                <FilterSearch setSearchQuery={setSearchQuery}/>
                <ul className='list-disc list-inside text-center'>
                    {/* <li><a className='underline' href='#'>IEEE_CC</a></li>
                    <li><a className='underline' href='#'>Google_Toros</a></li>
                    <li><a className='underline' href='#'>Dr_Izaddoost_Club</a></li> */}
                    {organizations.map((org, index) => (
                        <li key={index}><a className='underline' href='#'>{org["display_name"]}</a></li>
                    ))}
                </ul>
            </div>
        </div>
    );
}
