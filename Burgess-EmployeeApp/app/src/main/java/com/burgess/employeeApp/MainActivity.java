package com.burgess.employeeApp;

import android.support.v7.app.ActionBarActivity;
import android.app.Activity;
import android.bluetooth.BluetoothAdapter;
import android.bluetooth.BluetoothManager;
import android.net.wifi.WifiManager;
import android.net.ConnectivityManager;
import android.content.Context;
import android.content.Intent;
import android.os.Bundle;
import android.view.Menu;
import android.view.MenuItem;
import android.widget.Toast;

public class MainActivity extends ActionBarActivity {

	private final static int REQUEST_ENABLE_BT = 55;
	
	private BluetoothMetadataThread bluetoothThread;

	@Override
	protected void onCreate(Bundle savedInstanceState) {
		super.onCreate(savedInstanceState);
		setContentView(com.burgess.employeeApp.R.layout.activity_main);
		
		final BluetoothManager bluetoothManager = (BluetoothManager) getSystemService(Context.BLUETOOTH_SERVICE);
		BluetoothAdapter mBluetoothAdapter = bluetoothManager.getAdapter();

		if(mBluetoothAdapter == null)
		{
			Toast.makeText(getApplicationContext(), "Failed to get Bluetooth", Toast.LENGTH_LONG).show();
			return;
		}

		//asks user to enable Bluetooth for collection, if not already on
		if (!mBluetoothAdapter.isEnabled()) {
			Intent enableBtIntent = new Intent(BluetoothAdapter.ACTION_REQUEST_ENABLE);
			startActivityForResult(enableBtIntent, REQUEST_ENABLE_BT);
		}
		
		bluetoothThread = new BluetoothMetadataThread(bluetoothManager, (WifiManager) getSystemService(Context.WIFI_SERVICE), (ConnectivityManager) getSystemService(Activity.CONNECTIVITY_SERVICE), this);
		bluetoothThread.start();
	}

	@Override
	public boolean onCreateOptionsMenu(Menu menu) {
		// Inflate the menu; this adds items to the action bar if it is present.
		getMenuInflater().inflate(com.burgess.employeeApp.R.menu.main, menu);
		return true;
	}

	@Override
	public boolean onOptionsItemSelected(MenuItem item) {
		// Handle action bar item clicks here. The action bar will
		// automatically handle clicks on the Home/Up button, so long
		// as you specify a parent activity in AndroidManifest.xml.
		int id = item.getItemId();
		if (id == com.burgess.employeeApp.R.id.action_settings) {
			return true;
		}
		return super.onOptionsItemSelected(item);
	}
}
