import sys
import time
import pprint
import my_constant
import getopt

from web3.providers import HTTPProvider
from web3 import Web3

def wait_for_receipt(w3, tx_hash, poll_interval):
   while True:
       tx_receipt = w3.eth.getTransactionReceipt(tx_hash)
       if tx_receipt:
         return tx_receipt
       time.sleep(poll_interval)

if __name__ == '__main__':
  opts, args = getopt.getopt(sys.argv[1:], 'hi:p:c:w', ['help','wait','contract_address=','id=','http_provider='])

  # default http provider
  http_provider = 'http://localhost:8545'
  id,contract_address_str = 0,''
  wait = False

  for key, value in opts:
    # print('parameter',key,value)
    if key in ['-h', '--help']:
      print('upload design work info to ethereum')
      print('parameter:')
      print('-i\tdesign work id')
      print('-c\tcontract address')
      print('-w\twait for transaction end')
      print('-p\thttp provider address')
      sys.exit(0)
    if key in ['-i','--id']:
      id = int(value)
    if key in ['-p','--http_provider']:
      http_provider = value
    if key in ['-c','--contract_address']:
      contract_address_str = value
    if key in ['-w','--wait']:
      wait = True
  
  if id==0 or contract_address_str=='':
    print('empty parameter in id or contract address')
    sys.exit(-1)

  w3 = Web3(HTTPProvider(http_provider))
  connected = w3.isConnected()

  if not connected:
    print('not connected to',http_provider)
    sys.exit(-1)

  contract_address = w3.toChecksumAddress(contract_address_str)

  ca = w3.eth.contract(
     address=contract_address,
     abi=my_constant.CONTRACT_ABI)

  coin_base = w3.eth.accounts[0]
  print('eth.accounts[0]',coin_base)

  upload_address = ca.functions.showUploadAddress().call(get_block_id=False)
  print('upload address',upload_address)

  nt = int(time.time())

  # single create
  job = ca.functions.createDesignWork(id,nt)
  gas_estimate = job.estimateGas()
  print("single create gas estimate:",gas_estimate)

  # batch create
  id_arr = []
  time_arr = []
  for i in range(10):
    id_arr.append(i)
    time_arr.append(nt)
  batch_job = ca.functions.batchCreateDesignWork(id_arr,time_arr)
  batch_gas_estimate = batch_job.estimateGas()
  print("batch create gas estimate:",gas_estimate)

  tx_hash = job.transact({'gas':gas_estimate,'from':upload_address})
  batch_tx_hash = batch_job.transact({'gas':batch_gas_estimate,'from':upload_address})
  print("single transaction hash:",tx_hash.hex())
  print("batch transaction hash:",batch_tx_hash.hex())
  
  if wait:
    receipt = wait_for_receipt(w3, tx_hash, 1)
    print("single Transaction receipt mined!",'\n',receipt)

    rp = dict(receipt)
    print("single status:",rp['status'],"\n")

    batch_receipt = wait_for_receipt(w3, batch_tx_hash, 1)
    print("batch Transaction receipt mined!",'\n',receipt)

    batch_rp = dict(batch_receipt)
    print("batch status:",batch_rp['status'],"\n")