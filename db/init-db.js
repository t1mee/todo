const { MongoClient } = require('mongodb');

async function run() {
  try {
    const client = new MongoClient('mongodb://root:pass@localhost:5000', { useNewUrlParser: true });
    await client.connect();
    const db = client.db('todos');
    console.log('Database created!');
    await client.close();
  } catch (e) {
    console.error(e);
  }
}

run().catch(console.error);