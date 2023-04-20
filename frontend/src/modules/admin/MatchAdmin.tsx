import React from 'react';
import { useQuery } from 'react-query';

import MatchEditingTable from './match/MatchEditingTable';
import CreateMatchPanel from './match/CreateMatchPanel';
import Loading from '../shared/Loading';
import Error from '../shared/Error';
import { getFullMatches } from '@/lib/adminRequests';
import { QUERY_OPTIONS } from '@/lib/constants';
import styles from '@/styles/admin/MatchAdmin.module.css';
import type { MatchWithId } from '@/types/match';

interface MatchAdminProps {}

const MatchAdmin: React.FC<MatchAdminProps> = ({}) => {
  const { data, isLoading, isError } = useQuery(
    'allMatches',
    async () => getFullMatches(),
    QUERY_OPTIONS
  );

  if (isLoading) {
    return <Loading />
  };

  if (isError) {
    return <Error />
  }

  const Matches: MatchWithId[] = data || [];

  return (
    <div className={styles.container}>
      <MatchEditingTable inputMatchData={Matches} />
      <CreateMatchPanel />
    </div>
  );
};

export default MatchAdmin;