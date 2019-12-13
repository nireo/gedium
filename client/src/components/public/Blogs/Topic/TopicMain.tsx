import React from 'react';
import { Side } from './Side';
import { Main } from './Main';

type Props = {
  topic: string;
};

export const TopicMain: React.FC<Props> = ({ topic }) => {
  return (
    <div className="container">
      <div className="row" style={{ marginTop: '2rem' }}>
        <div className="col-md-9">
          <Main />
        </div>
        <div>
          <Side topic={topic} />
        </div>
      </div>
    </div>
  );
};
