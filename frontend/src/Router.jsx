import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Home } from './pages/Home';
import { Login } from './pages/Login';
import { LoadingScreen } from './components/LoadingScreen';
import { RequireAuth } from './components/RequireAuth';
import { AuthProvider } from './hooks';

export const Router = () => {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<LoadingScreen />} >
            <Route path='/login' element={<Login />} />
            <Route element={<RequireAuth />}>
              <Route path='/' element={<Home />} />
              <Route path='/foo' element={'Foo'} />
            </Route>
            <Route path='*' element={404} />
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
    )
};
