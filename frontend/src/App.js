import React from 'react';
import { Navigate } from 'react-router-dom';
import NavigationBar from './components/common/NavigationBar';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import { UserContext } from './contexts/userContext';

import PostPage from './pages/PostPage';

import Home from './pages/Home';
import Topics from './pages/Topics';
import Select from './pages/Select';
import Organizations from './pages/Organizations';

function App() {

    const [user, setUser] = React.useState({});
    const [loggedIn, setLoggedIn] = React.useState(false);

    React.useEffect(() => {
        const fetchData = async () => {
            try {
                const response = await fetch('http://localhost:3030/user/self', {
                    credentials: 'include'
                });
                if (response.status === 401) {
                    return setUser({ 'error': 'Not logged in' });
                }
                const data = await response.json();
                setUser(data);
                setLoggedIn(data['display_name'] !== undefined)
            } catch (error) {
                setUser({ 'error': error.message })
            }
        };

        fetchData();
    }, []);

    return (
      <BrowserRouter>
        <div className="app">
          <UserContext.Provider value={{ user, loggedIn }}>
            <NavigationBar />
            <Routes>
              <Route path='/' element={<Navigate to='/home' />} />
              <Route path='/home' element={<Home />} />
              <Route path='/posts/:postId' element={<PostPage />} />
              <Route path='/select' element={<Select />} />
              <Route path='/topics' element={<Topics />} />
              <Route path='/organizations' element={<Organizations />} />
              <Route path='*' element={<h1>Not Found</h1>} />
            </Routes>
          </UserContext.Provider>
        </div>
      </BrowserRouter>
    );
}

export default App;
